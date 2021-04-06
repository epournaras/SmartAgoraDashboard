//
//    COMMAND.GO -- Hive client command line interface
//
//    Copyright (c) 2016 Fragata Computer Systems AG
//    All rights reserved
//

package client

import (
    "os"
    "io"
    "io/ioutil"
    "bufio"
    "fmt"
    "encoding/json"
)

//
//    ---- Public constants
//

const (
    RoleAdmin = 0x01
    RoleUser  = 0x02
)

//
//    ---- CLI descriptor
//

type Command struct {
    client    *Client
    writer    io.Writer
    cookieDir string
    opMap     map[string] *opDesc
}

type opDesc struct {
    handler func(*Params, []string) error
    numArgs int
    usage   string
}

//
//    ---- Command dispatch
//

func NewCommand(client *Client, writer io.Writer, cookieDir string, roles int) *Command {
    c := &Command{client: client, writer: writer, cookieDir: cookieDir}
    c.initOpMap(roles)
    return c
}

func (c *Command) Do(params *Params, args []string) error {
    numArgs := len(args)
    if numArgs < 1 {
        return fmt.Errorf("Missing command")
    }
    op := args[0]
    desc, ok := c.opMap[op]
    if !ok {
        return fmt.Errorf("Invalid command '%s'", op)
    } 
    numArgs--;
    if numArgs != desc.numArgs {
        return fmt.Errorf("Usage: program %s %s", op, desc.usage) 
    }
    err := desc.handler(params, args[1:])
    if err != nil {
        return err
    }
    return nil
}

func (c *Command) initOpMap(roles int) {
    m := make(map[string] *opDesc)
    c.opMap = m;
    if roles & RoleAdmin != 0 {
        m["admin-root"] = &opDesc{c.AdminRoot, 0, ""}
        m["admin-setup"] = &opDesc{c.AdminSetup, 2, "reset-db setup-file"}
        m["admin-projects"] = &opDesc{c.AdminProjects, 0, "[params]"}
        m["admin-project"] = &opDesc{c.AdminProject, 1, "project-id"}
        m["admin-create-project"] = &opDesc{c.AdminCreateProject, 2, "project-id project-file"}
        m["admin-tasks"] = &opDesc{c.AdminTasks, 1, "[params] project-id"}
        m["admin-task"] = &opDesc{c.AdminTask, 2, "project-id task-id"}
        m["admin-create-tasks"] = &opDesc{c.AdminCreateTasks, 2, "project-id tasks-file"}
        m["admin-create-task"] = &opDesc{c.AdminCreateTask, 3, "project-id task-id task-file"}
        m["admin-complete-task"] = &opDesc{c.AdminCompleteTask, 2, "project-id task-id"}
        m["admin-disable-task"] = &opDesc{c.AdminDisableTask, 2, "project-id task-id"}
        m["admin-enable-task"] = &opDesc{c.AdminEnableTask, 2, "project-id task-id"}
        m["admin-assets"] = &opDesc{c.AdminAssets, 1, "[params] project-id"}
        m["admin-asset"] = &opDesc{c.AdminAsset, 2, "project-id asset-id"}
        m["admin-create-assets"] = &opDesc{c.AdminCreateAssets, 2, "project-id assets-data"}
        m["admin-users"] = &opDesc{c.AdminUsers, 1, "[params] project-id"}
        m["admin-user"] = &opDesc{c.AdminUser, 2, "project-id user-id"}
        m["admin-assignments"] = &opDesc{c.AdminAssignments, 1, "[params] project-id"}
    }
    if roles & RoleUser != 0 {
        m["project"] = &opDesc{c.Project, 1, "project-id"}
        m["tasks"] = &opDesc{c.Tasks, 1, "[params] project-id"}
        m["task"] = &opDesc{c.Task, 2, "project-id task-id"}
        m["asset"] = &opDesc{c.Asset, 2, "project-id asset-id"}
        m["user"] = &opDesc{c.User, 1, "project-id"}
        m["create-user"] = &opDesc{c.CreateUser, 2, "project-id user-file"}
        m["external-user"] = &opDesc{c.ExternalUser, 3, "project-id connect user-file"}
        m["favorites"] = &opDesc{c.Favorites, 1, "[params] project-id"}
        m["favorite"] = &opDesc{c.Favorite, 2, "project-id asset-id"}
        m["assignment"] = &opDesc{c.Assignment, 2, "project-id assignment-id"}
        m["assign-asset"] = &opDesc{c.AssignAsset, 3, "project-id task-id asset-id"}
        m["user-assignment"] = &opDesc{c.UserAssignment, 2, "project-id task-id"}
        m["user-create-assignment"] = &opDesc{c.UserCreateAssignment, 3, "project-id task-id assignment-file"}
    }
}

//
//    ---- Command handlers: administrator
//

func (c *Command) AdminRoot(params *Params, args []string) error {
    err := c.client.AdminRoot()
    if err != nil {
        return err
    }
    _, err = fmt.Fprintf(c.writer, "OK\n")
    return err
}

func (c *Command) AdminSetup(params *Params, args []string) error {
    project, tasks, assets, err := c.readSetup(args[1])
    if err != nil {
        return err
    }
    resetDb, err := parseBool(args[0])
    if err != nil {
        return err
    }
    projectId, numTasks, numAssets, err := c.client.AdminSetup(resetDb, project, tasks, assets)
    if err != nil {
        return err
    }
    _, err = fmt.Fprintf(c.writer, "Project ID: %s, tasks: %d, assets %d\n", projectId, numTasks, numAssets)
    return err
}

//
//    Projects
//

func (c *Command) AdminProjects(params *Params, args []string) error {
    projects, _, err := c.client.AdminProjects(params)
    if err != nil {
        return err
    }
    return c.writeProjects(projects)
}

func (c *Command) AdminProject(params *Params, args []string) error {
    project, err := c.client.AdminProject(args[0])
    if err != nil {
        return err
    }
    return c.writeProject(project)
}

func (c *Command) AdminCreateProject(params *Params, args []string) error {
    project, err := c.readProject(args[1])
    if err != nil {
        return err
    }
    project, err = c.client.AdminCreateProject(args[0], project)
    if err != nil {
        return err
    }
    return c.writeProject(project)
}

//
//    Tasks
//

func (c *Command) AdminTasks(params *Params, args []string) error {
    tasks, _, err := c.client.AdminTasks(args[0], params)
    if err != nil {
        return err
    }
    return c.writeTasks(tasks)
}

func (c *Command) AdminTask(params *Params, args []string) error {
    task, err := c.client.AdminTask(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeTask(task)
}

func (c *Command) AdminCreateTasks(params *Params, args []string) error {
    tasks, err := c.readTasks(args[1])
    if err != nil {
        return err
    }
    tasks, _, err = c.client.AdminCreateTasks(args[0], tasks)
    if err != nil {
        return err
    }
    return c.writeTasks(tasks)
}

func (c *Command) AdminCreateTask(params *Params, args []string) error {
    task, err := c.readTask(args[2])
    if err != nil {
        return err
    }
    task, err = c.client.AdminCreateTask(args[0], args[1], task)
    if err != nil {
        return err
    }
    return c.writeTask(task)
}

func (c *Command) AdminCompleteTask(params *Params, args []string) error {
    assets, _, err := c.client.AdminCompleteTask(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeAssets(assets)
}

func (c *Command) AdminDisableTask(params *Params, args []string) error {
    task, err := c.client.AdminDisableTask(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeTask(task)
}

func (c *Command) AdminEnableTask(params *Params, args []string) error {
    task, err := c.client.AdminEnableTask(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeTask(task)
}

//
//    Assets
//

func (c *Command) AdminAssets(params *Params, args []string) error {
    assets, _, err := c.client.AdminAssets(args[0], params)
    if err != nil {
        return err
    }
    return c.writeAssets(assets)
}

func (c *Command) AdminAsset(params *Params, args []string) error {
    asset, err := c.client.AdminAsset(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeAsset(asset)
}

func (c *Command) AdminCreateAssets(params *Params, args []string) error {
    assets, err := c.readAssets(args[1])
    if err != nil {
        return err
    }
    assets, _, err = c.client.AdminCreateAssets(args[0], assets)
    if err != nil {
        return err
    }
    return c.writeAssets(assets)
}

//
//    Users
//

func (c *Command) AdminUsers(params *Params, args []string) error {
    users, _, err := c.client.AdminUsers(args[0], params)
    if err != nil {
        return err
    }
    return c.writeUsers(users)
}

func (c *Command) AdminUser(params *Params, args []string) error {
    user, err := c.client.AdminUser(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeUser(user)
}

//
//    Assignments
//

func (c *Command) AdminAssignments(params *Params, args []string) error {
    assignments, _, err := c.client.AdminAssignments(args[0], params)
    if err != nil {
        return err
    }
    return c.writeAssignments(assignments)
}

//
//    ---- Command handlers: user
//

//
//    Project
//

func (c *Command) Project(params *Params, args []string) error {
    project, err := c.client.Project(args[0])
    if err != nil {
        return err
    }
    return c.writeProject(project)
}

//
//    Tasks
//

func (c *Command) Tasks(params *Params, args []string) error {
    tasks, _, err := c.client.Tasks(args[0], params)
    if err != nil {
        return err
    }
    return c.writeTasks(tasks)
}

func (c *Command) Task(params *Params, args []string) error {
    task, err := c.client.Task(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeTask(task)
}

//
//    Assets
//

func (c *Command) Asset(params *Params, args []string) error {
    asset, err := c.client.Asset(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeAsset(asset)
}

//
//    Users
//

func (c *Command) User(params *Params, args []string) error {
    userId, err := c.currentUserId(args[0])    
    if err != nil {
        return err
    }
    user, err := c.client.User(args[0], userId)
    if err != nil {
        return err
    }
    return c.writeUser(user)
}

func (c *Command) CreateUser(params *Params, args []string) error {
    user, err := c.readUser(args[1])
    if err != nil {
        return err
    }
    user, err = c.client.CreateUser(args[0], user)
    if err != nil {
        return err
    }
    err = c.setCurrentUserId(args[0], user.Id)
    if err != nil {
        return err
    }
    return c.writeUser(user)
}

func (c *Command) ExternalUser(params *Params, args []string) error {
    user, err := c.readUser(args[2])
    if err != nil {
        return err
    }
    connect, err := parseBool(args[1])
    if err != nil {
        return err
    }
    user, err = c.client.ExternalUser(args[0], connect, user)
    if err != nil {
        return err
    }
    return c.writeUser(user)
}

func (c *Command) Favorites(params *Params, args []string) error {
    userId, err := c.currentUserId(args[0])    
    if err != nil {
        return err
    }
    userFavorites, _, err := c.client.Favorites(args[0], userId, params) 
    if err != nil {
        return err
    }
    return c.writeUserFavorites(userFavorites)
}

func (c *Command) Favorite(params *Params, args []string) error {
    userId, err := c.currentUserId(args[0])    
    if err != nil {
        return err
    }
    assetId, action, err := c.client.Favorite(args[0], args[1], userId)
    if err != nil {
        return err
    }
    _, err = fmt.Fprintf(c.writer, "Asset ID: %s, action %s\n", assetId, action)
    return err
}

//
//    Assignments
//

func (c *Command) Assignment(params *Params, args []string) error {
    assignment, err := c.client.Assignment(args[0], args[1])
    if err != nil {
        return err
    }
    return c.writeAssignment(assignment)
}

func (c *Command) AssignAsset(params *Params, args []string) error {
    userId, err := c.currentUserId(args[0])    
    if err != nil {
        return err
    }
    assignment, err := c.client.AssignAsset(args[0], args[1], args[2], userId)
    if err != nil {
        return err
    }
    return c.writeAssignment(assignment)
}

func (c *Command) UserAssignment(params *Params, args []string) error {
    userId, err := c.currentUserId(args[0])    
    if err != nil {
        return err
    }
    assignment, err := c.client.UserAssignment(args[0], args[1], userId)
    if err != nil {
        return err
    }
    return c.writeAssignment(assignment)
}

func (c *Command) UserCreateAssignment(params *Params, args []string) error {
    userId, err := c.currentUserId(args[0])    
    if err != nil {
        return err
    }
    assignment, err := c.readAssignment(args[2])
    if err != nil {
        return err
    }
    assignment, err = c.client.UserCreateAssignment(args[0], args[1], userId, assignment)
    if err != nil {
        return err
    }
    return c.writeAssignment(assignment)
}

//
//    ---- Input functions
//

func (c *Command) readSetup(path string) (*Project, []Task, []Asset, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, nil, nil, err
    }
    var setup struct {
        Project *Project
        Tasks   []Task
        Assets  []Asset
    }
    err = json.Unmarshal(data, &setup)
    if err != nil {
        return nil, nil, nil, err
    }
    return setup.Project, setup.Tasks, setup.Assets, nil
}

func (c *Command) readProject(path string) (*Project, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    var project *Project
    err = json.Unmarshal(data, &project)
    if err != nil {
        return nil, err
    }
    return project, nil
}

func (c *Command) readTasks(path string) ([]Task, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    var tasks []Task
    err = json.Unmarshal(data, &tasks)
    if err != nil {
        return nil, err
    }
    return tasks, nil
}

func (c *Command) readTask(path string) (*Task, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    var task *Task
    err = json.Unmarshal(data, &task)
    if err != nil {
        return nil, err
    }
    return task, nil
}

func (c *Command) readAssets(path string) ([]Asset, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    var assets []Asset
    err = json.Unmarshal(data, &assets)
    if err != nil {
        return nil, err
    }
    return assets, nil
}

func (c *Command) readUser(path string) (*User, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    var user *User
    err = json.Unmarshal(data, &user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (c *Command) readAssignment(path string) (*Assignment, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    var assignment *Assignment
    err = json.Unmarshal(data, &assignment)
    if err != nil {
        return nil, err
    }
    return assignment, nil
}

//
//    ---- Output functions
//

func (c *Command) writeProjects(projects []Project) error {
    return c.write(projects)
}

func (c *Command) writeProject(project *Project) error {
    return c.write(project)
}

func (c *Command) writeTasks(tasks []Task) error {
    return c.write(tasks)
}

func (c *Command) writeTask(task *Task) error {
    return c.write(task)
}

func (c *Command) writeAssets(assets []Asset) error {
    return c.write(assets)
}

func (c *Command) writeAsset(asset *Asset) error {
    return c.write(asset)
}

func (c *Command) writeUsers(users []User) error {
    return c.write(users)
}

func (c *Command) writeUser(user *User) error {
    return c.write(user)
}

func (c *Command) writeAssignments(assignments []Assignment) error {
    return c.write(assignments)
}

func (c *Command) writeAssignment(assignment *Assignment) error {
    return c.write(assignment)
}

func (c *Command) writeUserFavorites(userFavorites UserFavorites) error {
    return c.write(userFavorites)
}

func (c *Command) write(v interface{}) error {
    data, err := json.MarshalIndent(v, "", "    ")
    if err != nil {
        return err
    }
    _, err = c.writer.Write(data)
    return err
}

//
//    ---- Utilities
//

func readFile(path string) ([]byte, error) {
    r, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer r.Close()
    data, err := ioutil.ReadAll(r)
    if err != nil {
        return nil, err
    }
    return data, nil
}

func (c *Command) setCurrentUserId(projectId, userId string) error {
    fn, err := c.cookiePath(projectId)
    if err != nil {
        return err
    }
    fp, err := os.Create(fn)
    if err != nil {
        return err
    }
    defer fp.Close()
    fmt.Fprintf(fp, "%s\n", userId)
    return nil
}

func (c *Command) currentUserId(projectId string) (string, error) {
    fn, err := c.cookiePath(projectId)
    if err != nil {
        return "", err
    }
    fp, err := os.Open(fn)
    if err != nil {
        return "", err
    }
    defer fp.Close()
    r := bufio.NewReader(fp)
    userId, err := r.ReadString('\n')
    if err != nil {
        return "", err
    }
    n := len(userId) - 1
    if n >= 0 && userId[n] == '\n' {
        userId = userId[0:n]
    }
    return userId, nil
}

func (c *Command) cookiePath(projectId string) (string, error) {
    if c.cookieDir == "" {
        return "", fmt.Errorf("Missing cookie directory")
    }
    return c.cookieDir + "/" + projectId + "_user_id.dat", nil
}

func parseBool(s string) (bool, error) {
    if s == "true" {
        return true, nil
    }
    if s == "false" {
        return false, nil
    }
    return false, fmt.Errorf("Invalid boolean argument '%s', must be 'true' or 'false'")
}

