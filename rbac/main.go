package rbac

type RoleType string

var (
	AccountViewer     RoleType = "Account Viewer"
	AccountManager    RoleType = "Account Manager"
	AccountAdmin      RoleType = "Account Admin"
	AccountView       string   = "account.view"
	AccountManage     string   = "account.manage"
	AccountDelete     string   = "account.delete"
	ProjectViewer     RoleType = "Project Viewer"
	ProjectManager    RoleType = "Project Manager"
	ProjectAdmin      RoleType = "Project Admin"
	ProjectView       string   = "Project.view"
	ProjectManage     string   = "Project.manage"
	ProjectDelete     string   = "Project.delete"
	TeamViewer        RoleType = "Team Viewer"
	TeamManager       RoleType = "Team Manager"
	TeamAdmin         RoleType = "Team Admin"
	TeamView          string   = "Team.view"
	TeamManage        string   = "Team.manage"
	TeamDelete        string   = "Team.delete"
	DeploymentViewer  RoleType = "Deployment Viewer"
	DeploymentManager RoleType = "Deployment Manager"
	DeploymentAdmin   RoleType = "Deployment Admin"
	DeploymentView    string   = "Deployment.view"
	DeploymentManage  string   = "Deployment.manage"
	DeploymentDelete  string   = "Deployment.delete"
	GuestViewer       RoleType = "Guest Viewer"
	GuestView         string   = "guest.view"
)

func GetRoles() []string {
	return []string{
		string(AccountViewer),
		string(AccountManager),
		string(AccountAdmin),
		string(ProjectViewer),
		string(ProjectManager),
		string(ProjectAdmin),
		string(TeamViewer),
		string(TeamManager),
		string(TeamAdmin),
		string(DeploymentViewer),
		string(DeploymentManager),
		string(DeploymentAdmin),
		string(GuestViewer),
	}
}

func GetRoleByString(r string) (RoleType, bool) {
	switch r {
	case "Account Viewer":
		return AccountViewer, true
	case "Account Manager":
		return AccountManager, true
	case "Account Admin":
		return AccountAdmin, true
	case "Project Viewer":
		return ProjectViewer, true
	case "Project Manager":
		return ProjectManager, true
	case "Project Admin":
		return ProjectAdmin, true
	case "Team Viewer":
		return TeamViewer, true
	case "Team Manager":
		return TeamManager, true
	case "Team Admin":
		return TeamAdmin, true
	case "Deployment Viewer":
		return DeploymentViewer, true
	case "Deployment Manager":
		return DeploymentManager, true
	case "Deployment Admin":
		return DeploymentAdmin, true
	case "Guest Viewer":
		return GuestViewer, true
	}

	return "", false
}

func GetPermissions(r string) []string {
	switch r {
	case "Account Viewer":
		return []string{AccountView}
	case "Account Manager":
		return []string{AccountView, AccountManage}
	case "Account Admin":
		return []string{AccountView, AccountManage, AccountDelete}
	case "Project Viewer":
		return []string{ProjectView}
	case "Project Manager":
		return []string{ProjectView, ProjectManage}
	case "Project Admin":
		return []string{ProjectView, ProjectManage, ProjectDelete}
	case "Team Viewer":
		return []string{TeamView}
	case "Team Manager":
		return []string{TeamView, TeamManage}
	case "Team Admin":
		return []string{TeamView, TeamManage, TeamDelete}
	case "Deployment Viewer":
		return []string{DeploymentView}
	case "Deployment Manager":
		return []string{DeploymentView, DeploymentManage}
	case "Deployment Admin":
		return []string{DeploymentView, DeploymentManage, DeploymentDelete}
	case "Guest Viewer":
		return []string{GuestView}
	}

	return []string{}
}
