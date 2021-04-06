
package hive

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

//
//    ---- HTTP server request multiplexer
//

// Starts up hive-server on the specified port, 
//     connecting to Elasticsearch at {esDomain}:{esPort} using the given index.
// Default parameters:
//     hive port: 8080
//     elasticsearch domain: localhost
//     elasticsearch port: 9200
//     elasticsearch index: hive

func (s *Server) Run() {
    log.Println("running hive-server on port", s.Port, "storing data in elasticsearch under index", s.Index)

    r := mux.NewRouter()
    r.StrictSlash(true)

    r.HandleFunc("/", s.RootHandler)

    //
    //    Administrator endpoints
    //

    // clear out db, configures elasticsearch and creates a project
    r.HandleFunc("/admin/setup", s.AdminSetupHandler)
    r.HandleFunc("/admin/setup/{DELETE_MY_DATABASE}", s.AdminSetupHandler)

    // return all projects in Hive
    r.HandleFunc("/admin/projects", s.AdminProjectsHandler).Methods("GET")

    // return project information
    r.HandleFunc("/admin/projects/{project_id}", s.AdminProjectHandler).Methods("GET")

    // create or update a project
    r.HandleFunc("/admin/projects/{project_id}", s.AdminCreateProjectHandler).Methods("POST")

    // return tasks in this project
    r.HandleFunc("/admin/projects/{project_id}/tasks", s.AdminTasksHandler).Methods("GET")

    // return task information
    r.HandleFunc("/admin/projects/{project_id}/tasks/{task_id}", s.AdminTaskHandler).Methods("GET")

    // import tasks into this project
    r.HandleFunc("/admin/projects/{project_id}/tasks", s.AdminCreateTasksHandler).Methods("POST")

    // create or update a task
    r.HandleFunc("/admin/projects/{project_id}/tasks/{task_id}", s.AdminCreateTaskHandler).Methods("POST")

    // enable and disable tasks
    r.HandleFunc("/admin/projects/{project_id}/tasks/{task_id}/enable", s.EnableTaskHandler).Methods("GET")
    r.HandleFunc("/admin/projects/{project_id}/tasks/{task_id}/disable", s.DisableTaskHandler).Methods("GET")

    // mark any assets completed for this task
    r.HandleFunc("/admin/projects/{project_id}/tasks/{task_id}/complete", s.CompleteTaskHandler)

    // return assets in this project
    r.HandleFunc("/admin/projects/{project_id}/assets", s.AdminAssetsHandler).Methods("GET")

    // get a single asset's data
    r.HandleFunc("/admin/projects/{project_id}/assets/{asset_id}", s.AdminAssetHandler)

    // import assets into this project
    r.HandleFunc("/admin/projects/{project_id}/assets", s.AdminCreateAssetsHandler).Methods("POST")

    // return users in this project
    r.HandleFunc("/admin/projects/{project_id}/users", s.AdminUsersHandler)

    // return a single user in this project
    r.HandleFunc("/admin/projects/{project_id}/users/{user_id}", s.AdminUserHandler)

    // return assignments in this project
    r.HandleFunc("/admin/projects/{project_id}/assignments", s.AdminAssignmentsHandler)

    //
    //    User endpoints
    //

    // return project information
    r.HandleFunc("/projects/{project_id}", s.ProjectHandler).Methods("GET")

    // return tasks in this project
    r.HandleFunc("/projects/{project_id}/tasks", s.TasksHandler).Methods("GET")

    // return task information
    r.HandleFunc("/projects/{project_id}/tasks/{task_id}", s.TaskHandler).Methods("GET")

    // return asset information
    r.HandleFunc("/projects/{project_id}/assets/{asset_id}", s.AssetHandler).Methods("GET")

    // return user information based on project session cookie
    r.HandleFunc("/projects/{project_id}/user", s.UserHandler).Methods("GET")

    // create a user based on json data posted
    r.HandleFunc("/projects/{project_id}/user", s.CreateUserHandler).Methods("POST")

    // look up user by external id, returns session token
    r.HandleFunc("/projects/{project_id}/user/external", s.ExternalUserHandler).Methods("POST")
    r.HandleFunc("/projects/{project_id}/user/external/{connect}", s.ExternalUserHandler).Methods("POST")

    // return a user's favorited ads
    r.HandleFunc("/projects/{project_id}/user/favorites", s.FavoritesHandler).Methods("GET")

    // favorite an asset
    r.HandleFunc("/projects/{project_id}/assets/{asset_id}/favorite", s.FavoriteHandler).Methods("GET")

    // return assignment information
    r.HandleFunc("/projects/{project_id}/assignments/{assignment_id}", s.AssignmentHandler).Methods("GET")

    // return a new assignment for task + asset + current user
    r.HandleFunc(
        "/projects/{project_id}/tasks/{task_id}/assets/{asset_id}/assignments", 
        s.AssignAssetHandler).Methods("GET")

    // return a new assignment for the given task + current user
    r.HandleFunc("/projects/{project_id}/tasks/{task_id}/assignments", s.UserAssignmentHandler).Methods("GET")

    // submit assignment (contribute, fill in form, etc)
    r.HandleFunc(
        "/projects/{project_id}/tasks/{task_id}/assignments", 
        s.UserCreateAssignmentHandler).Methods("POST")

    http.Handle("/", r)
    err := http.ListenAndServe(":"+s.Port, nil)
    if err != nil {
        log.Fatalf(err.Error())
    }
}

