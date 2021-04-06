
package hive

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "github.com/gorilla/mux"
    elastigo "github.com/jacqui/elastigo/lib"
)

//
//    ---- Private response data types
//

type projectResponse struct {
    Project Project
}

type projectsResponse struct {
    Projects []Project
    Meta     meta
}

type taskResponse struct {
    Task Task
}

type tasksResponse struct {
    Tasks []Task
    Meta  meta
}

type assetResponse struct {
    Asset Asset
}

type assetsResponse struct {
    Assets []Asset
    Meta   meta
}

type userResponse struct {
    User User
}

type usersResponse struct {
    Users []User
    Meta  meta
}

type assignmentResponse struct {
    Assignment Assignment
}
type assignmentsResponse struct {
    Assignments []Assignment
    Meta        meta
}

type favoriteResponse struct {
    AssetId string
    Action  string
}

type favoritesResponse struct {
    Favorites userFavorites
    Meta      meta
}

//
//    ---- Administrator request handlers
//

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
    endpointsJson := `{"status": "ok"}`
    s.wrapResponse(w, r, 200, []byte(endpointsJson))
}

//
//    Setup
//

// AdminSetupHandler -- clears out db, configures elasticsearch and creates a project
//     WARNING: this empties your database. Really.
// /admin/setup/{DELETE_MY_DATABASE} [post]
func (s *Server) AdminSetupHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    log.Println("Importing data into hive...")

    log.Println("Step 1: configuring elasticsearch.")
    indexExists, possible404 := s.EsConn.IndicesExists(s.Index)

    // for reasons mysterious to me, elastigo wraps all of the http pkg's functions
    // and does not check if the response to IndicesExists is a 404.
    // Elasticsearch will respond with a 404 if the index does not exist.
    // Here we check for this and correctly set the value of indexExists to false
    if possible404 != nil && possible404.Error() == "record not found" {
        indexExists = false

    } else if possible404 != nil {
        // otherwise some other error was thrown, so just 500 and give up here.
        s.wrapResponse(w, r, 500, s.wrapError(possible404))
        return
    }

    if vars["DELETE_MY_DATABASE"] == "YES_I_AM_SURE" && indexExists {
        // Delete existing hive index
        _, err := s.EsConn.DeleteIndex(s.Index)
        if err != nil {
            log.Println("Failed to delete index:", err)
            s.wrapResponse(w, r, 500, s.wrapError(err))
            return
        }
        log.Println("Deleted index", s.Index, ". I hope that was ok - you said you were sure!")
        indexExists = false
    } else if indexExists {
        giveUpErr := fmt.Errorf("Index '%s' exists. Use a different value or add 'YES_I_AM_SURE' to delete it: /admin/setup/YES_I_AM_SURE.", s.Index)
        s.wrapResponse(w, r, 500, s.wrapError(giveUpErr))
        return
    }

    if !indexExists {
        log.Println("Creating index", s.Index)
        // Create hive index
        _, err := s.EsConn.CreateIndex(s.Index)
        if err != nil {
            s.wrapResponse(w, r, 500, s.wrapError(err))
            return
        }
    }

    assignmentsBody := `{
        "assignments": {
            "properties": {
                "Asset": {
                    "properties": {
                        "Favorited": {
                            "type": "boolean"
                        },
                        "Id": {
                            "type": "string",
                            "index": "not_analyzed"
                        },
                        "Url": {
                            "type": "string",
                            "index": "not_analyzed"
                        }
                    }
                },
                "Id": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "Project": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "State": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "Task": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "User": {
                    "type": "string",
                    "index": "not_analyzed"
                }
            }
        }
    }`

    _, err := s.EsConn.DoCommand("PUT", 
        fmt.Sprintf("/%s/%s/_mapping", s.Index, "assignments"), nil, assignmentsBody)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    log.Println("Done configuring elasticsearch")

    log.Println("Step 2: creating project.")

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    var importedJson struct {
        Project Project
        Tasks   []Task
        Assets  []Asset
    }

    err = json.Unmarshal(body, &importedJson)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    s.ActiveProjectId = importedJson.Project.Id

    // store in elasticsearch
    _, err = s.EsConn.Index(s.Index, "projects", s.ActiveProjectId, nil, importedJson.Project)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    log.Println("Done creating project:", s.ActiveProjectId)

    log.Println("Step 3: importing tasks.")

    tasks, _, err := s.importTasks(importedJson.Tasks)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    log.Println("Done creating tasks:", len(tasks))

    log.Println("Step 4: adding assets.")

// -A.G. - Settings for "Projects" must be the same for all types
// (see ElasticSearch: the Definitive Guide [2.x] / Types and Mappings / Avoiding Type Gotchas); 
// hence "not_analyzed" is set here: 
/* TODO: Revise this
    assetsBody := `{
        "assets": {
            "properties": {
                "Id": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "Metadata": {
                    "properties": {
                        %s
                    }
                },
                "Project": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "SubmittedData": {
                    "type": "nested",
                    "include_in_parent": true,
                    "properties": {
                        %s
                    }
                },
                "Url": {
                    "type": "string"
                }
            }
        }
    }`
*/
/* 
  Type of SubmittedData changed from 'nested'to 'object'
  in conformance to Elasticsearch 2.x rules. It is actually not clear
  why it was declared 'nested' (is there any indexing by fields of arrays
  of objects involved?) 
*/
    assetsBody := `{
        "assets": {
            "properties": {
                "Id": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "Metadata": {
                    "properties": {
                        %s
                    }
                },
                "Project": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "SubmittedData": {
                    "type": "object",
                    "properties": {
                        %s
                    }
                },
                "Url": {
                    "type": "string"
                }
            }
        }
    }`

    project, err := s.FindProject(s.ActiveProjectId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    var metaProperties []string
    for _, metaProp := range project.MetaProperties {
        metaProperties = append(metaProperties, 
            fmt.Sprintf(`"%s": { "type": "%s", "index": "not_analyzed" }`, metaProp.Name, metaProp.Type))
    }
    metaPropertiesString := strings.Join(metaProperties, ",")

    var taskProperties []string
    for _, task := range tasks {
        taskProperties = append(taskProperties, fmt.Sprintf(`"%s": { "type": "object" }`, task.Name))
    }
    taskPropertiesString := strings.Join(taskProperties, ",")
    assetsMapping := fmt.Sprintf(assetsBody, metaPropertiesString, taskPropertiesString)

    _, err = s.EsConn.DoCommand("PUT", fmt.Sprintf("/%s/%s/_mapping", s.Index, "assets"), nil, assetsMapping)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    assets, err := s.importAssets(importedJson.Assets)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    log.Println("Done adding", len(assets), "assets")

    report := []byte(
        fmt.Sprintf(
            `{"status":"200 OK", "Project": "%s", "Tasks": "%d", "Assets": "%d"}`, 
                s.ActiveProjectId, len(tasks), len(assets)))
    s.wrapResponse(w, r, 200, report)
    return
}

//
//    Projects
//

// AdminProjectsHandler -- returns a paginated list of projects in Hive
// /admin/projects [get]
func (s *Server) AdminProjectsHandler(w http.ResponseWriter, r *http.Request) {
    queryParams := r.URL.Query()
    p := Params{
        From:    defaultQuery(queryParams, "from", "0"),
        Size:    defaultQuery(queryParams, "size", "10"),
        SortBy:  defaultQuery(queryParams, "sortBy", "Id"),
        SortDir: defaultQuery(queryParams, "sortDir", "asc"),
    }

    projects, m, err := s.FindProjects(p)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    resp := projectsResponse{
        Projects: projects,
        Meta:     m,
    }
    projectsJson, err := json.Marshal(resp)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, projectsJson)
}

// AdminProjectHandler -- returns a project by ID
// /admin/projects/{project_id} [get]
func (s *Server) AdminProjectHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    var project *Project
    var err error

    project, err = s.FindProject(s.ActiveProjectId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    resp := projectResponse{
        Project: *project,
    }
    projectJson, err := json.Marshal(resp)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, projectJson)
}

// AdminCreateProjectHandler -- creates or updates a project
// /admin/projects/{project_id} [post]
func (s *Server) AdminCreateProjectHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    var project *Project
    var err error

    project, err = s.CreateProject(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    resp := projectResponse{
        Project: *project,
    }
    projectJson, err := json.Marshal(resp)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, projectJson)
}

//
//    Tasks
//

// AdminTasksHandler -- returns a paginated tasks in a project
// /admin/projects/{project_id}/tasks [get]
func (s *Server) AdminTasksHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    queryParams := r.URL.Query()
    p := Params{
        From:    defaultQuery(queryParams, "from", "0"),
        Size:    defaultQuery(queryParams, "size", "10"),
        SortBy:  defaultQuery(queryParams, "sortBy", "Name"),
        SortDir: defaultQuery(queryParams, "sortDir", "asc"),
    }

    tasks, m, err := s.FindTasks(p)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    // format the json response
    tasksResponse := &tasksResponse{
        Tasks: tasks,
        Meta:  m,
    }
    tasksJson, err := json.Marshal(tasksResponse)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, tasksJson)
}

// AdminTaskHandler -- returns info for a single task by ID
// /admin/projects/{project_id}/tasks/{task_id} [get]
func (s *Server) AdminTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    taskId := vars["task_id"]
    if !strings.HasPrefix(vars["task_id"], s.ActiveProjectId) && vars["task_id"] != "" {
        taskId = s.ActiveProjectId + "-" + vars["task_id"]
    }

    task, err := s.FindTask(taskId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    taskJson, err := json.Marshal(taskResponse{
        Task: *task,
    })
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, taskJson)
}

// AdminCreateTasksHandler -- creates or updates tasks in a project
// /admin/projects/{project_id}/tasks [post]
func (s *Server) AdminCreateTasksHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    tasks, m, err := s.CreateTasks(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    tasksResponse := &tasksResponse{
        Tasks: tasks,
        Meta:  m,
    }
    tasksJson, err := json.Marshal(tasksResponse)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, tasksJson)
}

// AdminCreateTaskHandler -- creates or updates a task in a project
// /projects/{project_id}/tasks/{task_id} [get]
func (s *Server) AdminCreateTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    task, err := s.CreateTask(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    taskJson, err := json.Marshal(taskResponse{
        Task: *task,
    })
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, taskJson)
}

// DisableTaskHandler -- makes a task unavailable for assignment by disabling it
// /admin/projects/{project_id}/tasks/{task_id}/disable [get]
func (s *Server) DisableTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]
    taskId := vars["task_id"]
    taskName := taskId
    if !strings.HasPrefix(vars["task_id"], s.ActiveProjectId) && vars["task_id"] != "" {
        taskName = s.ActiveProjectId + "-" + taskName
    }

    task, err := s.UpdateTaskState(taskName, "waiting")
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    taskJson, err := json.Marshal(taskResponse{
        Task: *task,
    })

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, taskJson)
    return
}

// EnableTaskHandler -- makes a task available for assignment by enabling it
// /admin/projects/{project_id}/tasks/{task_id}/enable [get]
func (s *Server) EnableTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]
    taskId := vars["task_id"]
    taskName := taskId
    if !strings.HasPrefix(vars["task_id"], s.ActiveProjectId) && vars["task_id"] != "" {
        taskName = s.ActiveProjectId + "-" + taskName
    }

    task, err := s.UpdateTaskState(taskName, "available")
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    taskJson, err := json.Marshal(taskResponse{
        Task: *task,
    })

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, taskJson)
    return
}

// CompleteTaskHandler -- updates assets matching task CompletionCriteria with SubmittedData
// /admin/projects/{project_id}/tasks/{task_id}/complete [get]
func (s *Server) CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]
    taskId := vars["task_id"]

    assets, err := s.CompleteTask(taskId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    assetsJson, err := json.Marshal(assetsResponse{
        Assets: assets,
    })
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    s.wrapResponse(w, r, 200, assetsJson)
}

//
//    Assets
//

// AdminAssetsHandler -- returns a paginated list of assets in a project
// /admin/projects/{project_id}/assets [get]
func (s *Server) AdminAssetsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    var assets []Asset
    var m meta
    var err error

    queryParams := r.URL.Query()
    p := Params{
        From:    defaultQuery(queryParams, "from", "0"),
        Size:    defaultQuery(queryParams, "size", "10"),
        Task:    defaultQuery(queryParams, "task", ""),
        State:   defaultQuery(queryParams, "state", ""),
        SortBy:  defaultQuery(queryParams, "sortBy", "Id"),
        SortDir: defaultQuery(queryParams, "sortDir", "asc"),
    }

    if p.State == "completed" {
        assets, m, err = s.FindAssetsWithDataForTask(p)
        if err != nil {
            s.wrapResponse(w, r, 500, s.wrapError(err))
            return
        }
    }

    if p.State == "" {
        assets, m, err = s.FindAssets(p)
        if err != nil {
            s.wrapResponse(w, r, 500, s.wrapError(err))
            return
        }
    }

    var assetsWithCounts []Asset
    for _, asset := range assets {
        assetWithCounts, err := s.CalculateAssetCounts(asset)
        if err != nil {
            log.Println(err)
        }
        assetsWithCounts = append(assetsWithCounts, assetWithCounts)
    }

    assetsResponse := &assetsResponse{
        Assets: assetsWithCounts,
        Meta:   m,
    }
    assetsJson, err := json.Marshal(assetsResponse)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assetsJson)
}

// AdminAssetHandler -- retrieves a single project asset defined by an id
// /admin/projects/{project_id}/assets/{asset_id} [get]
func (s *Server) AdminAssetHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    assetId := vars["asset_id"]
    s.ActiveProjectId = vars["project_id"]

    asset, err := s.FindAsset(assetId)
    if err != nil {
        log.Println("failed finding asset", assetId, "because:", err)
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    assetWithCounts, err := s.CalculateAssetCounts(*asset)
    if err != nil {
        log.Println(err)
    }

    // format the json response
    resp := assetResponse{
        Asset: assetWithCounts,
    }

    assetJson, err := json.Marshal(resp)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assetJson)
}

// AdminCreateAssetsHandler -- creates assets in a project
// /admin/projects/{project_id}/assets [post]
func (s *Server) AdminCreateAssetsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    assets, err := s.CreateAssets(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    m := &meta{
        Total: len(assets),
        From:  0,
        Size:  10,
    }
    assetsJson, err := json.Marshal(&assetsResponse{
        Assets: assets,
        Meta:   *m,
    })
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assetsJson)
}

//
//    Users
//

// AdminUsersHandler -- returns a paginated list of users in a project
// /admin/projects/{project_id}/users [get]
func (s *Server) AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    queryParams := r.URL.Query()
    p := Params{
        From:     defaultQuery(queryParams, "from", "0"),
        Size:     defaultQuery(queryParams, "size", "10"),
        Task:     defaultQuery(queryParams, "task", ""),
        State:    defaultQuery(queryParams, "state", ""),
        SortBy:   defaultQuery(queryParams, "sortBy", "Id"),
        SortDir:  defaultQuery(queryParams, "sortDir", "asc"),
        Verified: defaultQuery(queryParams, "verified", ""),
    }

    _, err := s.EsConn.Refresh(s.Index)
    if err != nil {
        return
    }
    users, m, err := s.FindUsers(p)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    var assetIds []string
    assetQuery := `{ 
        "query": { 
            "query_string": { 
                "default_field": "Verified", 
                "query": "true" 
            } 
        }, 
        "aggs": { 
            "assets": { 
                "terms": { 
                    "field": "Id", 
                    "size": 0 
                } 
            } 
        } 
    }`
    assetResults, _ := s.EsConn.Search(s.Index, "assets", nil, assetQuery)
    var a assetAgg
    _ = json.Unmarshal(assetResults.Aggregations, &a)

    for _, b := range a.Assets.Buckets {
        assetIds = append(assetIds, b.Id)
    }
    assetIdString := "\"" + strings.Join(assetIds, "\", \"") + "\""
    for _, user := range users {
        if user.Counts["Assignments"] > 0 {
            verifyQuery := fmt.Sprintf(`{
                "query": {
                    "bool": {
                        "must": [
                        {
                            "terms": {
                                "assignments.Asset.Id": [%s]
                            }
                        },
                        {
                            "term": {
                                "assignments.User": "%s" 
                            } 
                        } 
                        ], 
                        "must_not": [ 
                        { 
                            "term": { 
                                "assignments.State": "skipped" 
                            } 
                        }, 
                        { 
                            "term": { 
                                "assignments.State": "unfinished" 
                            } 
                        } 
                        ] 
                    } 
                }, 
                "from": 0, 
                "size": %d
            }`, assetIdString, user.Id, user.Counts["Assignments"])
            verifyResults, _ := s.EsConn.Search(s.Index, "assignments", nil, verifyQuery)
            verifiedCount := verifyResults.Hits.Total
            user.Counts["VerifiedAssets"] = verifiedCount
            _, _ = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
        }
    }

    usersResponse := &usersResponse{
        Users: users,
        Meta:  m,
    }
    usersJson, err := json.Marshal(usersResponse)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, usersJson)
}

// AdminUserHandler -- returns a single user in a project by ID
// /admin/projects/{project_id}/users/{user_id} [get]
func (s *Server) AdminUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    _, err := s.EsConn.Refresh(s.Index)
    if err != nil {
        return
    }
    user, err := s.FindUser(vars["user_id"])
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    if user.Counts["Assignments"] > 0 {
        var assetIds []string
        assetQuery := `{ 
            "query": { 
                "query_string": { 
                    "default_field": "Verified", 
                    "query": "true" 
                } 
            }, 
            "aggs": { 
                "assets": { 
                    "terms": { 
                        "field": "Id", 
                        "size": 0 
                    } 
                } 
            } 
        }`
        assetResults, _ := s.EsConn.Search(s.Index, "assets", nil, assetQuery)
        var a assetAgg
        _ = json.Unmarshal(assetResults.Aggregations, &a)

        for _, b := range a.Assets.Buckets {
            assetIds = append(assetIds, b.Id)
        }
        assetIdString := "\"" + strings.Join(assetIds, "\", \"") + "\""
        verifyQuery := fmt.Sprintf(`{
            "query": {
                "bool": {
                    "must": [
                    {
                        "terms": {
                            "assignments.Asset.Id": [%s]
                        }
                    },
                    {
                        "term": {
                            "assignments.User": "%s" 
                        } 
                    } 
                    ], 
                    "must_not": [ 
                    { 
                        "term": { 
                            "assignments.State": "skipped" 
                        } 
                    }, 
                    { 
                        "term": { 
                            "assignments.State": "unfinished" 
                        } 
                    } 
                    ] 
                } 
            }, 
            "from": 0, 
            "size": %d
        }`, assetIdString, user.Id, user.Counts["Assignments"])
        verifyResults, _ := s.EsConn.Search(s.Index, "assignments", nil, verifyQuery)
        verifiedCount := verifyResults.Hits.Total
        user.Counts["VerifiedAssets"] = verifiedCount
        _, _ = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
    }

    userJson, err := json.Marshal(user)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    s.wrapResponse(w, r, 200, userJson)
}

//
//    Assignments
//

// AdminAssignmentsHandler -- returns a paginated list of assignments in a task
// /admin/projects/{project_id}/assignments [get]
func (s *Server) AdminAssignmentsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    queryParams := r.URL.Query()
    p := Params{
        From:    defaultQuery(queryParams, "from", "0"),
        Size:    defaultQuery(queryParams, "size", "10"),
        Task:    defaultQuery(queryParams, "task", ""),
        State:   defaultQuery(queryParams, "state", ""),
        SortBy:  defaultQuery(queryParams, "sortBy", "Id"),
        SortDir: defaultQuery(queryParams, "sortDir", "asc"),
    }

    assignments, m, err := s.FindAssignments(p)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    assignmentsResponse := &assignmentsResponse{
        Assignments: assignments,
        Meta:        m,
    }
    assignmentsJson, err := json.Marshal(assignmentsResponse)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assignmentsJson)
}

//
//    ---- User request handlers
//

//
//    Projects
//

// ProjectHandler -- returns a project by ID
// /projects/{project_id} [get]
func (s *Server) ProjectHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    var project *Project
    var err error

    project, err = s.FindProject(vars["project_id"])
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    // format the json response
    resp := projectResponse{
        Project: *project,
    }
    projectJson, err := json.Marshal(resp)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, projectJson)
}

//
//    Tasks
//

// TasksHandler -- returns a paginated tasks in a project
// /projects/{project_id}/tasks [get]
func (s *Server) TasksHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    queryParams := r.URL.Query()
    p := Params{
        From:    defaultQuery(queryParams, "from", "0"),
        Size:    defaultQuery(queryParams, "size", "10"),
        SortBy:  defaultQuery(queryParams, "sortBy", "Name"),
        SortDir: defaultQuery(queryParams, "sortDir", "asc"),
    }
    tasks, m, err := s.FindTasks(p)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    // format the json response
    tasksResponse := &tasksResponse{
        Tasks: tasks,
        Meta:  m,
    }
    tasksJson, err := json.Marshal(tasksResponse)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, tasksJson)
}

// TaskHandler -- returns public info for a single task by ID
// /projects/{project_id}/tasks/{task_id} [get]
func (s *Server) TaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    taskId := vars["task_id"]
    if !strings.HasPrefix(vars["task_id"], s.ActiveProjectId) && vars["task_id"] != "" {
        taskId = s.ActiveProjectId + "-" + vars["task_id"]
    }

    task, err := s.FindTask(taskId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    taskJson, err := json.Marshal(taskResponse{
        Task: *task,
    })
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, taskJson)
}

//
//    Assets
//

// AssetHandler -- returns public info for a single asset by ID
// /projects/{project_id}/assets/{asset_id} [get]
func (s *Server) AssetHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    assetId := vars["asset_id"]
    s.ActiveProjectId = vars["project_id"]

    asset, err := s.FindAsset(assetId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    resp := assetResponse{
        Asset: *asset,
    }
    assetJson, err := json.Marshal(resp)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assetJson)
}

// FavoritesHandler -- returns a paginated list of favorited assets for the current user
// /projects/{project_id}/user/favorites [get]
func (s *Server) FavoritesHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    sessionCookieName := s.ActiveProjectId + "_user_id"
    userId := s.FindCookieValue(r, sessionCookieName)
    user, err := s.FindUser(userId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    queryParams := r.URL.Query()
    p := Params{
        From: defaultQuery(queryParams, "from", "0"),
        Size: defaultQuery(queryParams, "size", "10"),
    }

    from, _ := strconv.Atoi(p.From)
    size, _ := strconv.Atoi(p.Size)

    m := meta{
        Total: len(user.Favorites),
        From:  from,
        Size:  size,
    }

    resp := favoritesResponse{
        Favorites: user.Favorites,
        Meta:      m,
    }
    favoritesJson, err := json.Marshal(resp)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, favoritesJson)
}

// FavoriteHandler -- toggles favoriting on an asset for the current user
// /projects/{project_id}/assets/{asset_id}/favorite [get]
func (s *Server) FavoriteHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    // find the asset
    asset, err := s.FindAsset(vars["asset_id"])
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    // find the user
    sessionCookieName := s.ActiveProjectId + "_user_id"
    userId := s.FindCookieValue(r, sessionCookieName)
    user, err := s.FindUser(userId)
    if user == nil {
        s.wrapResponse(w, r, 500, s.wrapError(errors.New("Favoriting assets requires a valid user.")))
        return
    }

    faveResponse := favoriteResponse{AssetId: asset.Id, Action: "favorited"}

    if len(asset.Counts) <= 0 {
        asset.Counts = Counts{
            "Favorites":   0,
            "Assignments": 0,
            "finished":    0,
            "skipped":     0,
            "unfinished":  0,
        }
    }
    if len(user.Favorites) <= 0 {
        user.Favorites = userFavorites{}
    }
    // is this asset in the user's favorites?
    _, ok := user.Favorites[asset.Id]
    if ok {
        delete(user.Favorites, asset.Id)
        faveResponse.Action = "unfavorited"
        if asset.Counts["Favorites"] > 0 {
            asset.Counts["Favorites"] -= 1
        }
    } else {
        // add the asset to the user's favorites
        user.Favorites[asset.Id] = *asset
        asset.Counts["Favorites"] += 1
    }
    user.Counts["Favorites"] = len(user.Favorites)

    _, err = s.EsConn.Index(s.Index, "assets", asset.Id, nil, asset)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    _, err = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    responseJson, err := json.Marshal(faveResponse)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    s.wrapResponse(w, r, 200, responseJson)
}

//
//    Users
//

// UserHandler -- returns info for the current user, creating a matching record if none found
// /projects/{project_id}/user [get]
func (s *Server) UserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    // user id is stored in a cookie named according to the project
    sessionCookieName := s.ActiveProjectId + "_user_id"

    // look for project's user session cookie
    userId := s.FindCookieValue(r, sessionCookieName)

    // try to find a matching user
    user, err := s.FindUser(userId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    // FindUser returns nil if no matching user is found
    if user == nil {
        tmpUser, err := s.CreateUserFromMissingCookieValue(userId)
        if err != nil {
            s.wrapResponse(w, r, 500, s.wrapError(err))
            return
        }
        user = &tmpUser
    }

    userJson, err := json.Marshal(user)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, userJson)
}

// CreateUserHandler -- creates a user in a project
// /projects/{project_id}/user [post]
func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    user, err := s.CreateUser(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
    }

    userJson, err := json.Marshal(user)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, userJson)
}

// ExternalUserHandler -- finds or creates a user by external ID
// /projects/{project_id}/user/external/{connect} [post]
func (s *Server) ExternalUserHandler(w http.ResponseWriter, r *http.Request) {
    var user *User
    var externalUser User
    var err error

    vars := mux.Vars(r)
    connectAccounts := vars["connect"]
    s.ActiveProjectId = vars["project_id"]

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    var lookupData struct {
        Id         string
        ExternalId string
    }

    err = json.Unmarshal(body, &lookupData)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    query := elastigo.Search(s.Index).Type("users").Filter(
        elastigo.Filter().Terms("ExternalId", lookupData.ExternalId),
        elastigo.Filter().Terms("Project", s.ActiveProjectId),
    )
    results, err := query.Result(&s.EsConn)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    resultCount := results.Hits.Total
    if resultCount != 0 {
        err = json.Unmarshal(*results.Hits.Hits[0].Source, &externalUser)
        if err != nil {
            s.wrapResponse(w, r, 500, s.wrapError(err))
            return
        }

        if externalUser.ExternalId == "0" {
            resultCount = 0
        }
    }

    // found no matching users
    if resultCount == 0 {
        userId := lookupData.Id

        if userId == "" && lookupData.Id != "" {
            userId = lookupData.Id
        }

        if userId == "" {
            // no ${project_id}_user_id set, create a new user
            tmpUser, err := s.CreateExternalUser(lookupData.ExternalId)
            if err != nil {
                s.wrapResponse(w, r, 500, s.wrapError(err))
                return
            }
            user = &tmpUser

        } else {
            // ${project_id}_user_id set, try looking up the user
            tmpUser, err := s.FindUser(userId)
            if err != nil {
                s.wrapResponse(w, r, 500, s.wrapError(err))
                return
            }

            user = tmpUser
            // found a user, set the externalId on it
            if user != nil {
                user.ExternalId = lookupData.ExternalId
                _, err = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
                if err != nil {
                    s.wrapResponse(w, r, 500, s.wrapError(err))
                    return
                }

            } else {
                // failed finding a user for that cookie (how would we get here?)
                *user, err = s.CreateExternalUser(lookupData.ExternalId)
                if err != nil {
                    s.wrapResponse(w, r, 500, s.wrapError(err))
                    return
                }
            }
        }
    }

    // found a matching user
    if resultCount == 1 {
        err = json.Unmarshal(*results.Hits.Hits[0].Source, &externalUser)
        if err != nil {
            s.wrapResponse(w, r, 500, s.wrapError(err))
            return
        }

        if connectAccounts == "" {
            user = &externalUser
        } else {
            userId := lookupData.Id
            tmpUser, err := s.FindUser(userId)
            if err != nil {
                s.wrapResponse(w, r, 500, s.wrapError(err))
                return
            }
            user = tmpUser
            if user != nil {
                user.ExternalId = lookupData.ExternalId

                // merge all the things

                // first: contribution counts
                for key, count := range externalUser.Counts {
                    user.Counts[key] += count
                }

                // second: favorites
                for key, value := range externalUser.Favorites {
                    user.Favorites[key] = value
                }

                user.Counts["VerifiedAssets"] = len(user.VerifiedAssets)

                _, err = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
                if err != nil {
                    s.wrapResponse(w, r, 500, s.wrapError(err))
                    return
                }

                // now, kill the other account
                var args map[string]interface{}
                _, err := s.EsConn.Delete(s.Index, "users", externalUser.Id, args)
                if err != nil {
                    s.wrapResponse(w, r, 500, s.wrapError(err))
                    return
                }
            }
        }
    }

    if resultCount > 1 {
        s.wrapResponse(w, r, 500, s.wrapError(errors.New("found more than one user with this externalId")))
        return
    }

    userJson, err := json.Marshal(user)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, userJson)
    return
}

//
//    Assignments
//

// AssignmentHandler -- returns public info for a single assignment by ID
// /projects/{project_id}/assignments/{assignment_id} [get]
func (s *Server) AssignmentHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]
    assignmentId := vars["assignment_id"]

    assignment, err := s.FindAssignment(assignmentId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    // format the json response
    resp := assignmentResponse{
        Assignment: *assignment,
    }
    assignmentJson, err := json.Marshal(resp)

    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assignmentJson)
}

// AssignAssetHandler -- finds or creates an unfinished assignment for the given asset, task and current user.
// /projects/{project_id}/tasks/{task_id}/assets/{asset_id}/assignments [get]
func (s *Server) AssignAssetHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]
    taskId := vars["task_id"]
    assetId := vars["asset_id"]

    // make sure taskId includes the active project
    if !strings.HasPrefix(vars["task_id"], s.ActiveProjectId) && vars["task_id"] != "" {
        taskId = s.ActiveProjectId + "-" + vars["task_id"]
    }

    // get user id from session cookie
    userId := s.FindCookieValue(r, s.ActiveProjectId+"_user_id")
    if userId == "" {
        userError := errors.New("Assignments can't be created without a user.")
        s.wrapResponse(w, r, 500, s.wrapError(userError))
        return
    }

    assignment, err := s.CreateAssetAssignment(taskId, userId, assetId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    assignJson, err := json.Marshal(assignment)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assignJson)
}

// UserAssignmentHandler -- finds or creates an unfinished task assignment for the current user.
// /projects/{project_id}/tasks/{task_id}/assignments [get]
func (s *Server) UserAssignmentHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]
    taskId := vars["task_id"]
    if !strings.HasPrefix(vars["task_id"], s.ActiveProjectId) && vars["task_id"] != "" {
        taskId = s.ActiveProjectId + "-" + vars["task_id"]
    }

    // get user id from session cookie
    sessionCookie, err := r.Cookie(s.ActiveProjectId + "_user_id")
    if err != nil { 
        // TODO: figure out how to avoid getting here;
        // frontend should check for user cookie before calling assign
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    userId := sessionCookie.Value

    assignment, err := s.CreateAssignment(taskId, userId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    assignJson, err := json.Marshal(assignment)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assignJson)
}

/* A.G. - 07.11.2016
// UserCreateAssignmentHandler -- finishes a task assignment & assigns a new one for the current user.
// /projects/{project_id}/tasks/{task_id}/assignments [post]
func (s *Server) UserCreateAssignmentHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]
    taskId := vars["task_id"]
    if !strings.HasPrefix(vars["task_id"], s.ActiveProjectId) && vars["task_id"] != "" {
        taskId = s.ActiveProjectId + "-" + vars["task_id"]
    }

    // get user id from session cookie
    userId := s.FindCookieValue(r, s.ActiveProjectId+"_user_id")

    _, err := s.UpdateAssignment(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    assignment, err := s.CreateAssignment(taskId, userId)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    assignJson, err := json.Marshal(assignment)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assignJson)
    return
}
*/

// UserCreateAssignmentHandler -- finishes a task assignment.
// /projects/{project_id}/tasks/{task_id}/assignments [post]
func (s *Server) UserCreateAssignmentHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    s.ActiveProjectId = vars["project_id"]

    assignment, err := s.UpdateAssignment(r.Body)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }

    assignJson, err := json.Marshal(assignment)
    if err != nil {
        s.wrapResponse(w, r, 500, s.wrapError(err))
        return
    }
    s.wrapResponse(w, r, 200, assignJson)
    return
}

//
//    ---- Handler utilities
//

// Looks for a cookie named 'cookieName' in the request.
// If the cookie is found, returns its value.
// Otherwise returns an empty string.
func (s *Server) FindCookieValue(r *http.Request, cookieName string) (cookieValue string) {
    cookie, err := r.Cookie(cookieName)

    // failed to find the cookie
    if err != nil {
        return ""
    }

    // cookie is empty
    if len(cookie.Value) == 0 || cookie.Value == "" {
        return ""
    }
    // found the cookie, return its value
    return cookie.Value
}

func defaultQuery(q url.Values, name string, defaultVal string) (val string) {
    qVal := q.Get(name)
    if qVal == "" {
        qVal = defaultVal
    }
    return qVal
}

// wrapResponse is a convenience function to consistently format responses with the right headers
func (s *Server) wrapResponse(w http.ResponseWriter, r *http.Request, statusCode int, data []byte) {

    w.Header().Set("Content-Type", "application/json")

    origin := r.Header.Get("Origin")
    if origin == "" {
        origin = r.Host
    }
    if origin != "" {
        w.Header().Set("Access-Control-Allow-Origin", origin)
    }

    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
    w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
    w.WriteHeader(statusCode)
    w.Write(data)
}

// wrapError is a convenience function to consistently format errors in json responses
func (s *Server) wrapError(err error) (formattedError []byte) {
    formattedError = []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error()))
    log.Println(string(formattedError))
    return formattedError
}

