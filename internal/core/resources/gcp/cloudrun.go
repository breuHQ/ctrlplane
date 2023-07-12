package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/iam/apiv1/iampb"
	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/shared"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

type (
	CloudRunConstructor struct {
	}

	CloudRun struct {
		ID                         gocql.UUID
		Cpu                        string
		Memory                     string
		Generation                 uint8
		Envs                       map[string]string
		OutputEnvs                 map[string]string
		Region                     string // from blueprint
		Image                      string // from workload
		Config                     string
		Name                       string
		Revision                   string
		LastRevision               string
		MinInstances               int32
		MaxInstances               int32
		AllowUnauthenticatedAccess bool
		Project                    string
		ServiceName                string
	}

	Workload struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	// AutoGenerated struct {
	// 	Properties struct {
	// 		Generation string `json:"generation"`
	// 		CPU        string `json:"cpu"`
	// 		Memory     string `json:"memory"`
	// 	} `json:"properties"`
	// 	Output struct {
	// 		Env []struct {
	// 			URL string `json:"url"`
	// 		} `json:"env"`
	// 	} `json:"output"`
	// }

	// config struct {
	// 	properties struct {
	// 		generation string `json:generation`
	// 		cpu        string `json:cpu`
	// 		memory     string `json:memory`
	// 	}
	// 	output struct {
	// 		env map[string]string
	// 	}
	// }

	GCPConfig struct {
		Project string
	}
)

// Create creates cloud run resource
func (c *CloudRunConstructor) Create(name string, region string, config string, providerConfig string) (core.CloudResource, error) {
	cr := &CloudRun{Name: name, Region: region, Config: config}
	cr.AllowUnauthenticatedAccess = true
	cr.Cpu = "2000m"
	cr.Memory = "1024Mi"
	cr.MinInstances = 1
	cr.MaxInstances = 8
	cr.Generation = 2

	// get gcp project from configuration
	pconfig := new(GCPConfig)
	err := json.Unmarshal([]byte(providerConfig), pconfig)
	if err != nil {
		shared.Logger().Error("Unable to parse provider config for cloudrun")
		return nil, err
	}

	cr.Project = pconfig.Project

	shared.Logger().Info("cloud run", "object", providerConfig, "umarshaled", pconfig, "project", cr.Project)
	return cr, nil
}

// CreateFromJson creates a CloudRun object from JSON
func (c *CloudRunConstructor) CreateFromJson(data []byte) core.CloudResource {
	cr := &CloudRun{}
	json.Unmarshal(data, cr)
	return cr
}

// Marshal marshals the CloudRun object
func (r *CloudRun) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Provision provisions the cloud resource
func (r *CloudRun) Provision(ctx workflow.Context) (workflow.Future, error) {

	// do nothing, the infra will be provisioned with deployment
	return nil, nil
}

// DeProvision deprovisions the cloudrun resource
func (r *CloudRun) DeProvision() error {

	return nil
}

// UpdateTraffic updates the traffic distribution on latest and previous revision as per the input
// parameter trafficpcnt is the percentage traffic to be deployed on latest revision
func (r *CloudRun) UpdateTraffic(ctx workflow.Context, trafficpcnt int32) error {

	// UpdateTraffic will execute a workflow to update the resource. This workflow is not directly called
	// from provisioninfra workflow to avoid passing resource interface as argument
	w := &Workflows{}

	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("CloudRun"),
			shared.WithWorkflowBlockID(r.Name),
			shared.WithWorkflowElement("UpdateTraffic"),
		)

	cctx := workflow.WithChildOptions(ctx, opts)

	shared.Logger().Info("Executing Update traffic workflow")

	err := workflow.
		ExecuteChildWorkflow(cctx, w.UpdateTraffic, r, trafficpcnt).Get(cctx, nil)

	if err != nil {
		shared.Logger().Error("Could not execute UpdateTraffic workflow", "error", err)
		return err
	}
	return nil
}

func (r *CloudRun) Deploy(ctx workflow.Context, wl []core.Workload) error {
	shared.Logger().Info("deploying", "cloudrun", r, "workload", wl)

	if len(wl) != 1 {
		shared.Logger().Error("Cannot deploy more than one workloads on cloud run", "number of workloads", len(wl))
		return errors.New("multiple workloads defined for cloud run")
	}

	// provision with execute a workflow to provision the resources. This workflow is not directly called
	// from provisioninfra workflow to avoid passing resource interface as argument
	w := &Workflows{}

	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("CloudRun"),
			shared.WithWorkflowBlockID(r.Name),
			shared.WithWorkflowElement("Deploy"),
		)

	cctx := workflow.WithChildOptions(ctx, opts)

	shared.Logger().Info("starting DeployCloudRun workflow")

	err := workflow.
		ExecuteChildWorkflow(cctx, w.DeployCloudRun, r, wl[0]).Get(cctx, r)

	if err != nil {
		shared.Logger().Error("Could not start DeployCloudRun workflow", "error", err)
		return err
	}
	return nil
}

func (w *Workflows) DeployCloudRun(ctx workflow.Context, r *CloudRun, wl *core.Workload) (*CloudRun, error) {

	r.ServiceName = wl.Name
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)
	err := workflow.ExecuteActivity(actx, activities.GetNextRevision, r).Get(actx, r)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: GetNextRevision", "error", err)
		return r, err
	}

	err = workflow.ExecuteActivity(actx, activities.DeployRevision, r, wl).Get(actx, nil)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: DeployDummy", "error", err)
		return r, err
	}
	return r, nil
}

// UpdateTraffic workflow executes UpdateTrafficActivity
func (w *Workflows) UpdateTraffic(ctx workflow.Context, r *CloudRun, trafficpcnt int32) error {

	shared.Logger().Info("Distributing traffic between revisions", r.Revision, r.LastRevision)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)
	err := workflow.ExecuteActivity(actx, activities.UpdateTrafficActivity, r, trafficpcnt).Get(ctx, r)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: UpdateTrafficActivity", "error", err)
		return err
	}
	return nil
}

// GetNextRevision Gets next revision Name to be deployed
// TODO: save the active resource's data on each deployment and on next deployment trigger get the associated data from the saved deployment.
func (a *Activities) GetNextRevision(ctx context.Context, r *CloudRun) (*CloudRun, error) {
	revision := r.GetFirstRevision()
	r.LastRevision = ""

	// get the deployed service, if not found then it will be first revision
	svc := r.GetService(ctx)
	if svc != nil {
		rev := svc.Template.Revision
		r.LastRevision = rev

		// revision name would be <service name>-<revision number> e.g first revision for helloworld service would be helloworld-0, second will be helloworld-1
		ss := strings.Split(rev, r.ServiceName+"-")
		revVersion, _ := strconv.Atoi(ss[1])
		revVersion++
		revision = r.ServiceName + "-" + strconv.Itoa(revVersion)
	}
	r.Revision = revision
	activity.GetLogger(ctx).Info("Next revision", "name", revision)
	return r, nil
}

// DeployRevision deploys a new revision on CloudRun if the service is already created.
// If no service is running, then it will create a new service and deploy first revision
func (a *Activities) DeployRevision(ctx context.Context, r *CloudRun, wl *core.Workload) error {

	logger := activity.GetLogger(ctx)
	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		logger.Error("could not create service client", "error", err)
	}

	defer client.Close()

	// Create service if this is the first revision
	if r.Revision == r.GetFirstRevision() {
		service := r.GetServiceTemplate(ctx, wl)
		logger.Info("deploying service", "data", service, "parent", r.GetParent(), "ID", wl.Name)
		csr := &runpb.CreateServiceRequest{Parent: r.GetParent(), Service: service, ServiceId: wl.Name}
		op, err := client.CreateService(ctx, csr)

		if err != nil {
			logger.Error("Could not create service", "Error", err)
			return err
		}

		logger.Info("waiting for service creation")
		op.Wait(ctx)
		// otherwise create a new revision and route 50% traffic to it
	} else {
		req := &runpb.GetServiceRequest{Name: r.GetParent() + "/services/" + wl.Name}
		service, err := client.GetService(ctx, req)

		logger.Info("50 percent traffic to latest", "revision", r.Revision)
		tt := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: 50}
		tt1 := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION, Revision: r.LastRevision, Percent: 50}
		service.Traffic = []*runpb.TrafficTarget{tt, tt1}

		if err != nil {
			logger.Error("could not get service", "Error", err)
			return err
		}

		service.Template.Revision = r.Revision
		usr := &runpb.UpdateServiceRequest{Service: service}
		op, err := client.UpdateService(ctx, usr)
		if err != nil {
			logger.Error("could not update service", "Error", err)
			return err
		}

		logger.Info("waiting for service revision update")
		op.Wait(ctx)
	}

	// Allow access to all users
	if r.AllowUnauthenticatedAccess {
		r.AllowAccessToAll(ctx)
	}

	return nil
}

// UpdateTrafficActivity updates traffic percentage on a cloud run resource
// This cannot be done in the workflow because of the blocking updateservice call
func (a *Activities) UpdateTrafficActivity(ctx context.Context, r *CloudRun, trafficpcnt int32) error {

	logger := activity.GetLogger(ctx)
	logger.Info("Update traffic", "revision", r.Revision, "percentage", trafficpcnt)
	service := r.GetService(ctx)
	cntxt := context.Background()
	serviceClient, err := run.NewServicesRESTClient(cntxt)
	if err != nil {
		logger.Error("New service rest client", "Error", err)
		return nil
	}
	defer serviceClient.Close()

	if r.Revision == r.GetFirstRevision() {
		ttc := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: 100}
		service.Traffic = []*runpb.TrafficTarget{ttc}
	} else {
		ttc := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: trafficpcnt}
		ttp := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION, Revision: r.LastRevision, Percent: 100 - trafficpcnt}
		service.Traffic = []*runpb.TrafficTarget{ttc, ttp}
	}

	req := &runpb.UpdateServiceRequest{Service: service}
	lro, err := serviceClient.UpdateService(cntxt, req)

	if err != nil {
		logger.Error("Update service", "Error", err)
		return err
	} else {
		logger.Info("waiting for service update")
		lro.Wait(cntxt)
	}
	return nil
}

func (r *CloudRun) GetServiceClient() (*run.ServicesClient, error) {
	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		shared.Logger().Error("New service rest client", "error", err)
		return nil, err
	}

	return client, err
}

// GetService gets a cloud run service from GCP
func (r *CloudRun) GetService(ctx context.Context) *runpb.Service {

	logger := activity.GetLogger(ctx)
	serviceClient, err := run.NewServicesRESTClient(ctx)
	if err != nil {
		shared.Logger().Error("New service rest client", "Error", err)
		return nil
	}
	defer serviceClient.Close()

	svcpath := r.GetParent() + "/services/" + r.ServiceName
	req := &runpb.GetServiceRequest{Name: svcpath}
	svc, err := serviceClient.GetService(ctx, req)

	if err != nil {
		logger.Error("Get Service", "Error", err)
		return nil
	}

	return svc
}

// AllowAccessToAll Sets IAM policy to allow access to all users
func (r *CloudRun) AllowAccessToAll(ctx context.Context) error {

	logger := activity.GetLogger(ctx)
	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		logger.Error("New service rest client", "Error", err)
		return nil
	}
	defer client.Close()

	rsc := r.GetParent() + "/services/" + r.ServiceName
	binding := new(iampb.Binding)
	binding.Members = []string{"allUsers"}
	binding.Role = "roles/run.invoker"
	Iamreq := &iampb.SetIamPolicyRequest{Resource: rsc, Policy: &iampb.Policy{Bindings: []*iampb.Binding{binding}}}
	_, err = client.SetIamPolicy(context.Background(), Iamreq)
	if err != nil {
		logger.Error("Set policy", "Error", err)
		return err
	}
	return nil
}

// GetServiceTemplate creates and returns the revision template for cloud run from the workload to be deployed
// revision template specifies the resource requirements, image to be deployed and traffic distribution etc.
// this template will be used for first deployment only, from next deployments the already deployed template will be
// fetched from cloudrun and the same will be used for next revision
// TODO: the above design will not work if resource definition is changed
func (r *CloudRun) GetServiceTemplate(ctx context.Context, wl *core.Workload) *runpb.Service {

	activity.GetLogger(ctx).Info("setting service template for", "revision", r.Revision)
	resources := &runpb.ResourceRequirements{Limits: map[string]string{"cpu": r.Cpu, "memory": r.Memory}}

	// unmarshaling the container here assuming that container definition will be specific to a resource
	// this can be done at a common location if the container definition turns out to be same for all resources
	crworkload := &Workload{}
	json.Unmarshal([]byte(wl.Container), crworkload)

	container := &runpb.Container{Name: wl.Name, Image: crworkload.Image, Resources: resources}

	scaling := &runpb.RevisionScaling{MinInstanceCount: r.MinInstances, MaxInstanceCount: r.MaxInstances}

	rt := &runpb.RevisionTemplate{Containers: []*runpb.Container{container}, Scaling: scaling,
		ExecutionEnvironment: runpb.ExecutionEnvironment(r.Generation), Revision: r.Revision}

	service := &runpb.Service{Template: rt}

	tt := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: 100}
	service.Traffic = []*runpb.TrafficTarget{tt}

	return service
}

func (r *CloudRun) GetParent() string {
	return "projects/" + r.Project + "/locations/" + r.Region
}

func (r *CloudRun) GetFirstRevision() string {
	return r.ServiceName + "-0"
}
