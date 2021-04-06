//
//    CLIENT.GO -- Hive client API
//
//    Copyright (c) 2016 Fragata Computer Systems AG
//    All rights reserved
//

package client

import (
    "bytes"
    "strconv"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "net/http"
)

//
//    ---- Public data types
//

type Params struct {
    From     string
    Size     string
    SortBy   string
    SortDir  string
    Task     string
    State    string
    Verified string
}

type Meta struct {
    Total int
    From  int
    Size  int
}

type Project struct {
    Id              string
    Name            string
    Description     string
    AssetCount      int
    TaskCount       int
    UserCount       int
    AssignmentCount Counts
    MetaProperties  []MetaProperty
}

type Task struct {
    Id                 string
    Project            string
    Name               string
    Description        string
    CurrentState       string
    AssignmentCriteria AssignmentCriteria
    CompletionCriteria CompletionCriteria
}

type Asset struct {
    Id            string
    Project       string
    Url           string
    Name          string
    Metadata      map[string]interface{}
    SubmittedData SubmittedData
    Favorited     bool
    Verified      bool
    Counts        Counts
}

type User struct {
    Id             string
    Name           string
    Email          string
    Project        string
    ExternalId     string
    Counts         Counts
    Favorites      UserFavorites
    NewFavorites   UserFavorites
    VerifiedAssets []string
}

type Assignment struct {
    Id            string
    User          string
    Project       string
    Task          string
    Asset         Asset
    State         string
    SubmittedData SubmittedData
}

type Counts map[string]int

type MetaProperty struct {
    Name string
    Type string
}

type AssignmentCriteria struct {
    SubmittedData map[string]interface{}
}

type CompletionCriteria struct {
    Total    int
    Matching int
}

type SubmittedData map[string]interface{}

type UserFavorites map[string]Asset

//
//    ---- Client descriptor
//

type Client struct {
    scheme string
    host   string
    port   string
}

func NewClient(scheme, host, port string) *Client {
    return &Client{scheme: scheme, host: host, port: port}
}

//
//    ---- Admin interface
//

//
//    Setup
//

func (c *Client) AdminRoot() error {
    path := "/"
    return c.httpGet(path, nil, nil)
}

// AdminSetup -- finds or creates an unfinished task assignment for the current user
// /admin/setup/{DELETE_MY_DATABASE} [put]
func (c *Client) AdminSetup(
        resetDb bool, project *Project, tasks []Task, assets []Asset) (string, int, int, error) {
    var deleteMyDatabase string
    if resetDb {
        deleteMyDatabase = "/YES_I_AM_SURE"
    }
    path := fmt.Sprintf("/admin/setup%s", deleteMyDatabase)
    var data setupResponse
    err := c.httpPost(path, nil, &setupRequest{Project: project, Tasks: tasks, Assets: assets}, &data)
    if err != nil {
        return "", 0, 0, err
    }
    numTasks, err := strconv.Atoi(data.Tasks)
    if err != nil {
        return "", 0, 0, err
    }
    numAssets, err := strconv.Atoi(data.Assets)
    if err != nil {
        return "", 0, 0, err
    }
    return data.Project, numTasks, numAssets, nil
}

//
//    Projects
//

// AdminProjects -- returns a paginated list of projects in Hive
// /admin/projects [get]
func (c *Client) AdminProjects(params *Params) ([]Project, *Meta, error) {
    path := "/admin/projects"
    var data projectsResponse
    err := c.httpGet(path, params, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Projects, data.Meta, nil
    }

// AdminProject -- returns a project by ID
// /admin/projects/{project_id} [get]
func (c *Client) AdminProject(projectId string) (*Project, error) {
    path := fmt.Sprintf("/admin/projects/%s", projectId)
    var data projectResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Project, nil
}

// AdminCreateProject -- creates or updates a project
// /admin/projects/{project_id} [post]
func (c *Client) AdminCreateProject(projectId string, project *Project) (*Project, error) {
    path := fmt.Sprintf("/admin/projects/%s", projectId)
    var data projectResponse
    err := c.httpPost(path, nil, project, &data)
    if err != nil {
        return nil, err
    }
    return data.Project, nil
}

//
//    Tasks
//

// AdminTasks -- returns a paginated list of tasks in a project
// /admin/projects/{project_id}/tasks [get]
func (c *Client) AdminTasks(projectId string, params *Params) ([]Task, *Meta, error) {
    path := fmt.Sprintf("/admin/projects/%s/tasks", projectId)
    var data tasksResponse
    err := c.httpGet(path, params, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Tasks, data.Meta, nil
}

// AdminTask -- returns info for a single task by ID
// /admin/projects/{project_id}/tasks/{task_id} [get]
func (c *Client) AdminTask(projectId, taskId string) (*Task, error) {
    path := fmt.Sprintf("/admin/projects/%s/tasks/%s", projectId, taskId)
    var data taskResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Task, nil
}

// AdminCreateTasks -- creates or updates tasks in a project
// /admin/projects/{project_id}/tasks [post]
func (c *Client) AdminCreateTasks(projectId string, tasks []Task) ([]Task, *Meta, error) {
    path := fmt.Sprintf("/admin/projects/%s/tasks", projectId)
    var data tasksResponse
    err := c.httpPost(path, nil, &tasksRequest{Tasks: tasks}, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Tasks, data.Meta, nil
}

// AdminCreateTask -- creates or updates a task in a project
// /admin/projects/{project_id}/tasks/{task_id} [post]
func (c *Client) AdminCreateTask(projectId, taskId  string, task *Task) (*Task, error) {
    // ACHTUNG: Server ignores 'taskId'; 
    //     it derives 'task.Id' from 'projectId' and 'task.Name'
    path := fmt.Sprintf("/admin/projects/%s/tasks/%s", projectId, taskId)
    var data taskResponse
    err := c.httpPost(path, nil, task, &data)
    if err != nil {
        return nil, err
    }
    return data.Task, nil
}

// AdminCompleteTask -- updates assets matching task CompletionCriteria with SubmittedData
// /admin/projects/{project_id}/tasks/{task_id}/complete [get]
func (c *Client) AdminCompleteTask(projectId, taskId string) ([]Asset, *Meta, error) {
    path := fmt.Sprintf("/admin/projects/%s/tasks/%s/complete", projectId, taskId)
    var data assetsResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Assets, data.Meta, nil
}

// AdminDisableTask -- makes a task unavailable for assignment by disabling it
// /admin/projects/{project_id}/tasks/{task_id}/disable [get]
func (c *Client) AdminDisableTask(projectId, taskId string) (*Task, error) {
    path := fmt.Sprintf("/admin/projects/%s/tasks/%s/disable", projectId, taskId)
    var data taskResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Task, nil
}

// AdminEnableTask -- makes a task available for assignment by enabling it
// /admin/projects/{project_id}/tasks/{task_id}/enable [get]
func (c *Client) AdminEnableTask(projectId, taskId string) (*Task, error) {
    path := fmt.Sprintf("/admin/projects/%s/tasks/%s/enable", projectId, taskId)
    var data taskResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Task, nil
}

//
//    Assets
//

// AdminAssets -- returns a paginated list of assets in a project
// /admin/projects/{project_id}/assets [get]
func (c *Client) AdminAssets(projectId string, params *Params) ([]Asset, *Meta, error) {
    path := fmt.Sprintf("/admin/projects/%s/assets", projectId)
    var data assetsResponse
    err := c.httpGet(path, params, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Assets, data.Meta, nil
}

// AdminAsset -- retrieves a single project asset defined by an id
// /admin/projects/{project_id}/assets/{asset_id} [get]
func (c *Client) AdminAsset(projectId, assetId string) (*Asset, error) {
    path := fmt.Sprintf("/admin/projects/%s/assets/%s", projectId, assetId)
    var data assetResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Asset, nil
}

// AdminCreateAssets -- creates assets in a project
// /admin/projects/{project_id}/assets [post]
func (c *Client) AdminCreateAssets(projectId string, assets []Asset) ([]Asset, *Meta, error) {
    path := fmt.Sprintf("/admin/projects/%s/assets", projectId)
    var data assetsResponse
    err := c.httpPost(path, nil, &assetsRequest{Assets: assets}, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Assets, data.Meta, nil
}

//
//    Users
//

// AdminUsers -- returns a paginated list of users in a project
// /admin/projects/{project_id}/users [get]
func (c *Client) AdminUsers(projectId string, params *Params) ([]User, *Meta, error) {
    path := fmt.Sprintf("/admin/projects/%s/users", projectId)
    var data usersResponse
    err := c.httpGet(path, params, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Users, data.Meta, nil
}

// AdminUser -- returns a single user in a project by ID
// /admin/projects/{project_id}/users/{user_id} [get]
func (c *Client) AdminUser(projectId, userId string) (*User, error) {
    path := fmt.Sprintf("/admin/projects/%s/users/%s", projectId, userId)
/*    TODO: Revise this
    var data userResponse
*/
    var data *User
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
/*    TODO: Revise this
    return data.User, nil
*/
    return data, nil
} 

//
//    Assignments
//

// AdminAssignments -- returns a paginated list of assignments in a task
// /admin/projects/{project_id}/assignments [get]
func (c *Client) AdminAssignments(projectId string, params *Params) ([]Assignment, *Meta, error) {
    path := fmt.Sprintf("/admin/projects/%s/assignments", projectId)
    var data assignmentsResponse
    err := c.httpGet(path, params, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Assignments, data.Meta, nil
}

//
//    ---- User interface
//

//
//    Projects
//

// Project -- returns a project by ID
// /projects/{project_id} [get]
func (c *Client) Project(projectId string) (*Project, error) {
    path := fmt.Sprintf("/projects/%s", projectId)
    var data projectResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Project, nil
}

//
//    Tasks
//

// Tasks -- returns a paginated list of tasks in a project
// /projects/{project_id}/tasks [get]
func (c *Client) Tasks(projectId string, params *Params) ([]Task, *Meta, error) {
    path := fmt.Sprintf("/projects/%s/tasks", projectId)
    var data tasksResponse
    err := c.httpGet(path, params, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Tasks, data.Meta, nil
}

// Task -- returns public info for a single task by ID
// /projects/{project_id}/tasks/{task_id} [get]
func (c *Client) Task(projectId, taskId string) (*Task, error) {
    path := fmt.Sprintf("/projects/%s/tasks/%s", projectId, taskId)
    var data taskResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Task, nil
}

//
//    Assets
//

// Asset -- returns public info for a single asset by ID
// /projects/{project_id}/assets/{asset_id} [get]
func (c *Client) Asset(projectId, assetId string) (*Asset, error) {
    path := fmt.Sprintf("/projects/%s/assets/%s", projectId, assetId)
    var data assetResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Asset, nil
}

//
//    Users
//

// User -- returns info for the current user, creating a matching record if none found
// /projects/{project_id}/user [get]
func (c *Client) User(projectId, userId string) (*User, error) {
    path := fmt.Sprintf("/projects/%s/user", projectId)
/*    TODO: Revise this
    var data userResponse
*/
    var data *User
    err := c.httpGetWithUser(path, projectId, userId, &data)
    if err != nil {
        return nil, err
    }
/*    TODO: Revise this
    return data.User, nil
*/
    return data, nil
}

// CreateUser -- creates a user in a project
// /projects/{project_id}/user [post]
func (c *Client) CreateUser(projectId string, user *User) (*User, error) {
    path := fmt.Sprintf("/projects/%s/user", projectId)
/*    TODO: Revise this
    var data userResponse
*/
    var data *User
    err := c.httpPost(path, nil, user, &data)
    if err != nil {
        return nil, err
    }
/*    TODO: Revise this
    return data.User, nil
*/
    return data, nil
}

// ExternalUser -- finds or creates a user by external ID
// /projects/{project_id}/user/external/{connect} [post]
func (c *Client) ExternalUser(projectId string, connect bool, user *User) (*User, error) {
    path := fmt.Sprintf("/projects/%s/user/external", projectId)
    if connect {
        path += "/connect"
    }
/*    TODO: Revise this
    var data userResponse
*/
    var data *User
    err := c.httpPost(path, nil, user, &data)
    if err != nil {
        return nil, err
    }
/*    TODO: Revise this
    return data.User, nil
*/
    return data, nil
}

// Favorites -- returns a paginated list of favorited assets for the current user
// /projects/{project_id}/user/favorites [get]
func (c *Client) Favorites(projectId, userId string, params *Params) (UserFavorites, *Meta, error) {
    path := fmt.Sprintf("/projects/%s/user/favorites", projectId)
    var data favoritesResponse
    err := c.httpGetWithUser(path, projectId, userId, &data)
    if err != nil {
        return nil, nil, err
    }
    return data.Favorites, data.Meta, nil
}

// Favorite -- toggles favoriting on an asset for the current user
// /projects/{project_id}/assets/{asset_id}/favorite [get]
func (c *Client) Favorite(projectId, assetId, userId string) (string, string, error) {
    path := fmt.Sprintf("/projects/%s/assets/%s/favorite", projectId, assetId)
    var data favoriteResponse
    err := c.httpGetWithUser(path, projectId, userId, &data)
    if err != nil {
        return "", "", err
    }
    return data.AssetId, data.Action, nil
}

//
//    Assignments
//

// Assignment -- returns public info for a single assignment by ID
// /projects/{project_id}/assignments/{assignment_id} [get]
func (c *Client) Assignment(projectId, assignmentId string) (*Assignment, error) {
    path := fmt.Sprintf("/projects/%s/assignments/%s", projectId, assignmentId)
    var data assignmentResponse
    err := c.httpGet(path, nil, &data)
    if err != nil {
        return nil, err
    }
    return data.Assignment, nil
}

// AssignAsset -- finds or creates an unfinished assignment for the given asset, task and current user
// /projects/{project_id}/tasks/{task_id}/assets/{asset_id}/assignments [get]
func (c *Client) AssignAsset(projectId, taskId, assetId, userId string) (*Assignment, error) {
    path := fmt.Sprintf("/projects/%s/tasks/%s/assets/%s/assignments", projectId, taskId, assetId)
/* TODO: Revise this
    var data assignmentResponse
*/
    var data *Assignment
    err := c.httpGetWithUser(path, projectId, userId, &data)
    if err != nil {
        return nil, err
    }
/* TODO: Revise this
    return data.Assignment, nil
*/
    return data, nil
}

// UserAssignment -- finds or creates an unfinished task assignment for the current user
// /projects/{project_id}/tasks/{task_id}/assignments [get]
func (c *Client) UserAssignment(projectId, taskId, userId string) (*Assignment, error) {
    path := fmt.Sprintf("/projects/%s/tasks/%s/assignments", projectId, taskId)
/* TODO: Revise this
    var data assignmentResponse
*/
    var data *Assignment
    err := c.httpGetWithUser(path, projectId, userId, &data)
    if err != nil {
        return nil, err
    }
/* TODO: Revise this
    return data.Assignment, nil
*/
    return data, nil
}

// UserCreateAssignment -- finishes a task assignment & assigns a new one for the current user
// /projects/{project_id}/tasks/{task_id}/assignments [post]
func (c *Client) UserCreateAssignment(
        projectId, taskId, userId string, assignment *Assignment) (*Assignment, error) {
    path := fmt.Sprintf("/projects/%s/tasks/%s/assignments", projectId, taskId)
/* TODO: Revise this
    var data assignmentResponse
*/
    var data *Assignment
    err := c.httpPostWithUser(path, projectId, userId, assignment, &data)
    if err != nil {
        return nil, err
    }
/*
    return data.Assignment, nil
*/
    return data, nil
}

//
//    ---- Private data types
//

// NOTE: Fields are upper-case to facilitate default bindings at JSON ummarshaling

// requests

type setupRequest struct {
    Project *Project
    Tasks []Task
    Assets []Asset
}

type tasksRequest struct {
    Tasks []Task
}

type assetsRequest struct {
    Assets []Asset
}

// responses

type setupResponse struct {
    Project string
    Tasks   string
    Assets  string
}

type projectResponse struct {
    Project *Project
}

type projectsResponse struct {
    Projects []Project
    Meta     *Meta
}

type taskResponse struct {
    Task *Task
}

type tasksResponse struct {
    Tasks []Task
    Meta  *Meta
}

type assetResponse struct {
    Asset *Asset
}

type assetsResponse struct {
    Assets []Asset
    Meta   *Meta
}

/*    TODO: Revise this
type userResponse struct {
    User *User
}
*/

type usersResponse struct {
    Users []User
    Meta  *Meta
}

type favoriteResponse struct {
    AssetId string
    Action  string
}

type favoritesResponse struct {
    Favorites UserFavorites
    Meta      *Meta
}

type assignmentResponse struct {
    Assignment *Assignment
}

type assignmentsResponse struct {
    Assignments []Assignment
    Meta        *Meta
}

//
//    ---- HTTP interface
//

func (c *Client) httpGet(path string, params *Params, out interface{}) error {
    url := c.formatURL(path, params)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    // TODO: Handle status code 500 nicely
    //     (extract error message from JSON {"error":"message"}
    if resp.StatusCode != 200 {
        return httpError(resp.Status, body)
    }
    if out == nil {
        return nil
    }
    err = json.Unmarshal(body, out)
    if err != nil {
        return err
    }
    return nil
}

func (c *Client) httpGetWithUser(path string, projectId, userId string, out interface{}) error {
    url := c.formatURL(path, nil)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }
    req.AddCookie(&http.Cookie{Name: projectId+"_user_id", Value: userId})
    clnt := &http.Client{}
    resp, err := clnt.Do(req) 
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    // TODO: Handle status code 500 nicely
    //     (extract error message from JSON {"error":"message"}
    if resp.StatusCode != 200 {
        return httpError(resp.Status, body)
    }
    if out == nil {
        return nil
    }
    err = json.Unmarshal(body, out)
    if err != nil {
        return err
    }
    return nil
}

func (c *Client) httpPost(path string, params *Params, in, out interface{}) error {
    url := c.formatURL(path, params)
    body, err := json.Marshal(in)
    if err != nil {
        return err
    }
    resp, err := http.Post(url, "application/json", bytes.NewReader(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    // TODO: Handle status code 500 nicely
    //     (extract error message from JSON {"error":"message"}
    if resp.StatusCode != 200 {
        return httpError(resp.Status, body)
    }
    if out == nil {
        return nil
    }
    err = json.Unmarshal(body, out)
    if err != nil {
        return err
    }
    return nil
}

func (c *Client) httpPostWithUser(path string, projectId, userId string, in, out interface{}) error {
    url := c.formatURL(path, nil)
    body, err := json.Marshal(in)
    req, err := http.NewRequest("POST", url, bytes.NewReader(body))
    if err != nil {
        return err
    }
    req.AddCookie(&http.Cookie{Name: projectId+"_user_id", Value: userId})
    req.Header.Add("content-type", "application/json")
    clnt := &http.Client{}
    resp, err := clnt.Do(req) 
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    // TODO: Handle status code 500 nicely
    //     (extract error message from JSON {"error":"message"}
    if resp.StatusCode != 200 {
        return httpError(resp.Status, body)
    }
    if out == nil {
        return nil
    }
    err = json.Unmarshal(body, out)
    if err != nil {
        return err
    }
    return nil
}

//
//    ---- Request formatting
//

func (c *Client) formatURL(path string, params *Params) string {
    var port string
    if len(c.port) > 0 {
        port = ":" + c.port 
    }
    return c.scheme + "://" + c.host + port + path + formatQuery(params)  
}

func formatQuery(params *Params) string {
    if params == nil {
        return ""
    }
    query := ""
    query = appendParam(query, "from", params.From)
    query = appendParam(query, "size", params.Size)
    query = appendParam(query, "sortBy", params.SortBy)
    query = appendParam(query, "sortDir", params.SortDir)
    query = appendParam(query, "task", params.Task)
    query = appendParam(query, "state", params.State)
    query = appendParam(query, "verified", params.Verified)
    return query
}

func appendParam(query, name, value string) string {
    if (len(value) > 0) {
        if (len(query) == 0) {
            query = fmt.Sprintf("?%s=%s", name, value)
        } else {
            query += fmt.Sprintf("&%s=%s", name, value)
        }
    }
    return query
}

//
//    ---- Utility functions
//

func httpError(status string, body []byte) error {
    return fmt.Errorf("[%s]\n%s", status, body)
}

