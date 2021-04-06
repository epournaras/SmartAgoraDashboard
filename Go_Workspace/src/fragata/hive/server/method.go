
package hive

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "math/rand"
    "reflect"
    "strconv"
    "strings"
    elastigo "github.com/jacqui/elastigo/lib"
)

//
//    ---- Public data types
//

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
    Favorites      userFavorites
    NewFavorites   userFavorites
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

type userFavorites map[string]Asset

type Params struct {
    From     string
    Size     string
    SortBy   string
    SortDir  string
    Task     string
    State    string
    Verified string
}

type meta struct {
    Total int
    From  int
    Size  int
}

//
//    ----- Server descriptor
//

type Server struct {
    Port            string
    Index           string
    EsConn          elastigo.Conn
    ActiveProjectId string
}

func NewServer() *Server {
    return &Server{}
}

//
//    ----- Private Elasticsearch data types
//

type assetBucket struct {
    Id    string      `json:"key"`
    Count int         `json:"doc_count"`
    Users userBuckets `json:"users"`
}

type assetBuckets struct {
    Buckets []assetBucket `json:"buckets"`
}

type assetAgg struct {
    Assets assetBuckets `json:"assets"`
}

type userBucket struct {
    Id    string `json:"key"`
    Count int    `json:"doc_count"`
}

type userBuckets struct {
    Buckets []userBucket `json:"buckets"`
}

type assignmentBucket struct {
    Key   string `json:"key"`
    Count int    `json:"doc_count"`
}

type assignmentBuckets struct {
    Buckets []assignmentBucket `json:"buckets"`
}

type assignmentAgg struct {
    Assignments assignmentBuckets `json:"assignments"`
}

//
//    ---- Methods
//

//
//    Projects
//

// FindProjects returns all projects, tallying counts of assets, users, tasks and assignments for each.
func (s *Server) FindProjects(p Params) (projects []Project, m meta, err error) {
    query := elastigo.Search(s.Index).Type("projects").From(p.From).Size(p.Size)
    results, err := query.Result(&s.EsConn)

    if err != nil {
        return
    }

    resultCount := results.Hits.Total

    m.Total = resultCount
    m.From, _ = strconv.Atoi(p.From)
    m.Size, _ = strconv.Atoi(p.Size)
    if resultCount <= 0 {
        err = errors.New("No projects found")
        return

    } else {
        for _, hit := range results.Hits.Hits {
            var project Project
            rawMessage := hit.Source
            err = json.Unmarshal(*rawMessage, &project)
            if err != nil {
                return
            }

// -A.G. (Server.Count requires ActiveProjectId)
            s.ActiveProjectId = project.Id

            project.AssetCount, _ = s.Count("assets")
            project.UserCount, _ = s.Count("users")
            project.TaskCount, _ = s.Count("tasks")
            project.AssignmentCount, _ = s.CountAssignments()

            projects = append(projects, project)
        }
    }
    return
}

// FindProject looks up a project by id, tallying counts of assets, users, tasks and assignments.
func (s *Server) FindProject(id string) (project *Project, err error) {
    err = s.EsConn.GetSource(s.Index, "projects", id, nil, &project)
    if err != nil {
        return nil, err
    }
    project.AssetCount, _ = s.Count("assets")
    project.UserCount, _ = s.Count("users")
    project.TaskCount, _ = s.Count("tasks")

    project.AssignmentCount, _ = s.CountAssignments()

    return project, nil
}

// Creates or updates a project by parsing the JSON body of the request.
func (s *Server) CreateProject(requestBody io.Reader) (project *Project, err error) {
    body, err := ioutil.ReadAll(requestBody)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(body, &project)
    if err != nil {
        return nil, err
    }

    // store in elasticsearch
    _, err = s.EsConn.Index(s.Index, "projects", project.Id, nil, project)
    if err != nil {
        return nil, err
    }
    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return nil, err
    }

    return project, nil
}

// Count composes a simple elasticsearch query scoping results to the current project, 
//     returning a total of 'countWhat'
// This method is used to tally number of tasks and assets for instance.
func (s *Server) Count(countWhat string) (count int, err error) {
    var args map[string]interface{}

    projectQuery := fmt.Sprintf(`{ "query": { "term" : {"Project": "%s" } } }`, s.ActiveProjectId)
    countResponse, err := s.EsConn.Count(s.Index, countWhat, args, projectQuery)
    if err != nil {
        return
    }
    count = countResponse.Count
    return
}

//
//    Tasks
//

// FindTasks returns an array of tasks for the current project
func (s *Server) FindTasks(p Params) (tasks []Task, m meta, err error) {
    query := elastigo.Search(s.Index).Type("tasks").Filter(
        elastigo.Filter().Terms("Project", s.ActiveProjectId),
    ).From(p.From).Size(p.Size)
    if p.SortDir == "desc" {
        query = query.Sort(
            elastigo.Sort(p.SortBy).Desc(),
        )
    } else {
        query = query.Sort(
            elastigo.Sort(p.SortBy).Asc(),
        )
    }
    results, err := query.Result(&s.EsConn)

    if err != nil {
        tasks = make([]Task, 0)
        return
    }

    for _, hit := range results.Hits.Hits {
        var task Task
        rawMessage := hit.Source
        err = json.Unmarshal(*rawMessage, &task)
        if err != nil {
            return
        }
        tasks = append(tasks, task)
    }
    return
}

// FindTask looks up a task by id
func (s *Server) FindTask(id string) (task *Task, err error) {
    err = s.EsConn.GetSource(s.Index, "tasks", id, nil, &task)
    if err != nil {
        return nil, err
    }
    return task, nil
}

// CreateTasks reads the request body POST'd to hive's admin to create/update tasks
func (s *Server) CreateTasks(requestBody io.Reader) (tasks []Task, m meta, err error) {
    body, err := ioutil.ReadAll(requestBody)
    if err != nil {
        return
    }

    var importedJson struct {
        Tasks []Task
    }
    err = json.Unmarshal(body, &importedJson)
    if err != nil {
        return
    }

    tasks, m, err = s.importTasks(importedJson.Tasks)
    if err != nil {
        return
    }

    return tasks, m, nil
}

// importTasks is a helper method called by CreateTasks that formats
//     the request body appropriately for saving tasks.
func (s *Server) importTasks(newTasks []Task) (tasks []Task, m meta, err error) {
    for _, task := range newTasks {
        if len(task.Name) == 0 {
            err = errors.New("Sorry, all tasks must specify a name.")
            return
        }
        task.Project = s.ActiveProjectId

        task.Id = strings.Join([]string{s.ActiveProjectId, strings.ToLower(task.Name)}, "-")
        if task.AssignmentCriteria.SubmittedData == nil {
            task.AssignmentCriteria.SubmittedData = make(map[string]interface{})
        }

        // store in elasticsearch, which will generate a unique id
        _, err := s.EsConn.Index(s.Index, "tasks", task.Id, nil, task)
        if err != nil {
            return tasks, m, err
        }
        tasks = append(tasks, task)
    }
    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return
    }

    m.Total = len(tasks)
    m.From = 0
    m.Size = len(tasks)

    return tasks, m, nil
}

// Creates or updates a task by parsing the JSON body of the request.
func (s *Server) CreateTask(requestBody io.Reader) (task *Task, err error) {
    body, err := ioutil.ReadAll(requestBody)
    if err != nil {
        return
    }

    err = json.Unmarshal(body, &task)
    if err != nil {
        return
    }

// -A.G. (copied from Server.importTasks)
    if len(task.Name) == 0 {
        err = errors.New("Sorry, a task must specify a name.")
        return
    }
    task.Project = s.ActiveProjectId

    task.Id = strings.Join([]string{s.ActiveProjectId, strings.ToLower(task.Name)}, "-")
    if task.AssignmentCriteria.SubmittedData == nil {
        task.AssignmentCriteria.SubmittedData = make(map[string]interface{})
    }
    _, err = s.EsConn.Index(s.Index, "tasks", task.Id, nil, task)
    if err != nil {
        return
    }

    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return
    }

    return task, nil
}

// UpdateTaskState is called from disable and enable TaskHandlers
// It sets the current state of a task (available, waiting)
func (s *Server) UpdateTaskState(taskId string, state string) (task *Task, err error) {
    task, err = s.FindTask(taskId)
    if err != nil {
        return nil, err
    }
    task.CurrentState = state
    _, err = s.EsConn.Index(s.Index, "tasks", task.Id, nil, task)
    if err != nil {
        return nil, err
    }
    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return nil, err
    }
    return
}

type SubmittedDataTracker struct {
    Value SubmittedData
    Count int
}

// CompleteTask uses the task's CompletionCriteria to find eligible assets for verification.
func (s *Server) CompleteTask(taskId string) ([]Asset, error) {
    var searchJson string
    var assets []Asset

    taskName := s.ActiveProjectId + "-" + taskId
    task, err := s.FindTask(taskName)
    if err != nil {
        return assets, err
    }

    query := `{
        "aggs": {
            "assets": {
                "terms": {
                    "field": "Asset.Id",
                    "size": 50000,
                    "min_doc_count": %d
                },
                "aggs": {
                    "users": {
                        "terms": {
                            "field": "User"
                        }
                    }
                }
            }
        },
        "query": {
            "filtered": {
                "filter": {
                    "bool": {
                        "must": [
                        {
                            "query": {
                                "match": {
                                    "assignments.Task": "%s"
                                }
                            }
                        },
                        {
                            "query": {
                                "match": {
                                    "Project": "%s"
                                }
                            }
                        },
                        {
                            "query": {
                                "match": {
                                    "State": "finished"
                                }
                            }
                        }
                        ]
                    }
                }
            }
        }
    }`

    searchJson = fmt.Sprintf(query, task.CompletionCriteria.Total, taskName, s.ActiveProjectId)
    log.Println(searchJson)

    results, err := s.EsConn.Search(s.Index, "assignments", nil, searchJson)
    if err != nil {
        return assets, err
    }

    log.Println("** Assignments count:", results.Hits.Total)
    var a assetAgg
    err = json.Unmarshal(results.Aggregations, &a)
    if err != nil {
        return nil, err
    }

    log.Println("** Assets Buckets:", len(a.Assets.Buckets))
    for _, b := range a.Assets.Buckets {
        if b.Count >= task.CompletionCriteria.Matching {
            log.Println("Completing asset", b.Id, "for task", task.Name)

            assignmentQuery := `{
                "query": {
                    "filtered": {
                        "filter": {
                            "bool": {
                                "must": [
                                {
                                    "query": {
                                        "match": {
                                            "Task": "%s"
                                        }
                                    }
                                },
                                {
                                    "query": {
                                        "match": {
                                            "Asset.Id": "%s"
                                        }
                                    }
                                },
                                {
                                    "query": {
                                        "match": {
                                            "Project": "%s"
                                        }
                                    }
                                },
                                {
                                    "query": {
                                        "match": {
                                            "State": "finished"
                                        }
                                    }
                                }
                                ]
                            }
                        }
                    }
                }
            }`
            assignmentSearchJson := fmt.Sprintf(assignmentQuery, taskName, b.Id, s.ActiveProjectId)
            log.Println(assignmentSearchJson)
            assignmentResults, err := s.EsConn.Search(s.Index, "assignments", nil, assignmentSearchJson)
            if err != nil {
                log.Println("error searching for matching assignment:", err)
                return nil, err
            }
            log.Println("** Matching assignments count:", assignmentResults.Hits.Total)

            var matchingAssignments []Assignment
            var sdTrackers []SubmittedDataTracker
            for _, assignmentHit := range assignmentResults.Hits.Hits {
                var matchingAssignment Assignment
                rawMessage := assignmentHit.Source
                err = json.Unmarshal(*rawMessage, &matchingAssignment)
                if err != nil {
                    log.Println(err)
                    continue
                }

                sdTrackers = collateSubmittedData(sdTrackers, matchingAssignment.SubmittedData)
                matchingAssignments = append(matchingAssignments, matchingAssignment)
            }

            log.Println("sdTrackers:", sdTrackers)
            for _, tracker := range sdTrackers {
                if tracker.Count >= task.CompletionCriteria.Matching {
                    log.Println("found", tracker.Count, "matching sds!")
                    asset, err := s.CompleteAsset(b.Id, *task, tracker.Value)
                    if err != nil {
                        log.Println("error completing asset", err)
                        continue
                    }
                    assets = append(assets, *asset)
                    for _, a := range matchingAssignments {
                        a.State = "verified"
                        log.Println("verifying assignment", a.Id)
                        _, err = s.EsConn.Index(s.Index, "assignments", a.Id, nil, a)
                        if err != nil {
                            log.Println("error saving assignment record:", err)
                        }
                    }
                    continue
                }
            }
        }
    }

    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return assets, err
    }

    return assets, err
}

func collateSubmittedData(sdt []SubmittedDataTracker, item SubmittedData) []SubmittedDataTracker {
    log.Println("---------------------------------------")
    log.Println("sdt size:", len(sdt))
    log.Println("sdt before:", sdt)
    log.Println("item:", item)
    foundIt := false
    for i, tracker := range sdt {
        if reflect.DeepEqual(tracker.Value, item) {
            log.Println("found a match")
            // we've seen this before
            tracker.Count += 1
            sdt[i] = tracker
            log.Println("count is now:", tracker.Count)
            foundIt = true
        }
    }
    log.Println("sdt after:", sdt)
    if !foundIt {
        log.Println("didn't find it")
        sdt = append(sdt, SubmittedDataTracker{
            Value: item,
            Count: 1,
        })
    }
    log.Println("---------------------------------------")
    return sdt
}

//
//    Assets
//

// FindAssets returns an array of assets in the current project, along with pagination meta information.
// 'from' and 'size' parameters determine the offset and limit passed to the database.
func (s *Server) FindAssets(p Params) (assets []Asset, m meta, err error) {
    query := elastigo.Search(s.Index).Type("assets").Filter(
        elastigo.Filter().Terms("Project", s.ActiveProjectId),
    ).From(p.From).Size(p.Size)
    if p.SortDir == "desc" {
        query = query.Sort(
            elastigo.Sort(p.SortBy).Desc(),
        )
    } else {
        query = query.Sort(
            elastigo.Sort(p.SortBy).Asc(),
        )
    }
    results, err := query.Result(&s.EsConn)

    if err != nil {
        return
    }

    resultCount := results.Hits.Total

    m.Total = resultCount
    m.From, _ = strconv.Atoi(p.From)
    m.Size, _ = strconv.Atoi(p.Size)

    for _, hit := range results.Hits.Hits {
        var asset Asset
        rawMessage := hit.Source
        err = json.Unmarshal(*rawMessage, &asset)
        if err != nil {
            return
        }
/*
        // use this when reindexing assets
        _, err = s.EsConn.Index(s.Index, "assets", asset.Id, nil, asset)
        if err != nil {
            return
        }
*/
        if len(asset.Counts) <= 0 {
            asset.Counts = Counts{
                "Favorites":   0,
                "Assignments": 0,
                "finished":    0,
                "skipped":     0,
                "unfinished":  0,
            }
        }
        assets = append(assets, asset)
    }
/*
    // use this when reindexing assets
    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return
    }
*/
    return
}

// FindAssetsWithDataForTask returns a list of assets in the current project and 
// given task with submitted/verified data, along with pagination meta information.
// 'from' and 'size' parameters determine the offset and limit passed to the database.
// 'sortBy' and 'sortDir' parameters determine ordering of results
func (s *Server) FindAssetsWithDataForTask(p Params) (assets []Asset, m meta, err error) {
    var exists []string

    taskParams := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }

    tasks, m, err := s.FindTasks(taskParams)
    if p.Task != "" {
        exists = append(exists, fmt.Sprintf(`{ "exists": { "field": "SubmittedData.%s" } }`, p.Task))
    } else {
        if err != nil {
            return
        }
        for _, t := range tasks {
            exists = append(exists, fmt.Sprintf(`{ "exists": { "field": "SubmittedData.%s" } }`, t.Name))
        }
    }
    searchQuery := `{
        "query": {
            "filtered": {
                "filter":   {
                    "bool": {
                        "must": [%s]
                    }
                }
            }
        },
        "from": %s,
        "size": %s,
        "sort": [ { "%s": { "order" : "%s" } } ]
    }`

    searchJson := fmt.Sprintf(searchQuery, strings.Join(exists, ", "), p.From, p.Size, p.SortBy, p.SortDir)
    log.Println(searchJson)
    results, err := s.EsConn.Search(s.Index, "assets", nil, searchJson)
    if err != nil {
        return
    }

    m.Total = results.Hits.Total
    m.From, _ = strconv.Atoi(p.From)
    m.Size, _ = strconv.Atoi(p.Size)

    for _, hit := range results.Hits.Hits {
        var asset Asset
        rawMessage := hit.Source
        err = json.Unmarshal(*rawMessage, &asset)
        if err != nil {
            return
        }
        assets = append(assets, asset)
    }
    return
}

// FindAsset looks up an asset by id.
func (s *Server) FindAsset(id string) (asset *Asset, err error) {
    err = s.EsConn.GetSource(s.Index, "assets", id, nil, &asset)
    if err != nil {
        return nil, err
    }
    return asset, nil
}

// Creates assets in this project by parsing the JSON body of the request.
func (s *Server) CreateAssets(requestBody io.Reader) (assets []Asset, err error) {
    body, err := ioutil.ReadAll(requestBody)
    if err != nil {
        return assets, err
    }

    var importedJson struct {
        Assets []Asset
    }
    err = json.Unmarshal(body, &importedJson)
    if err != nil {
        return assets, err
    }

    assets, err = s.importAssets(importedJson.Assets)
    if err != nil {
        return assets, err
    }
    return assets, nil

}

// importAssets is a helper method called by CreateAssets that formats the request body appropriately for saving assets.
func (s *Server) importAssets(newAssets []Asset) (assets []Asset, err error) {
    p := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }
    tasks, _, err := s.FindTasks(p)
    if err != nil {
        return assets, err
    }

    submittedData := SubmittedData{}
    for _, task := range tasks {
        submittedData[task.Name] = nil
    }

    for _, asset := range newAssets {
        if len(asset.Url) == 0 {
            return assets, errors.New("Sorry, all assets must specify a url.")
        }
        asset.Project = s.ActiveProjectId
        asset.SubmittedData = submittedData
        asset.Counts = Counts{
            "Favorites":   0,
            "Assignments": 0,
            "finished":    0,
            "skipped":     0,
            "unfinished":  0,
        }

        // store in elasticsearch, which will generate a unique id
        result, err := s.EsConn.Index(s.Index, "assets", "", nil, asset)
        if err != nil {
            return assets, err
        }

        // get the id, store it in the asset source in elasticsearch
        asset.Id = result.Id
        _, err = s.EsConn.Index(s.Index, "assets", asset.Id, nil, asset)
        if err != nil {
            return assets, err
        }

        if err == nil {
            assets = append(assets, asset)
        }
    }

    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return
    }

    return assets, nil
}

// CompleteAsset is called by CompleteTask to store verified submitted data on assets.
func (s *Server) CompleteAsset(assetId string, task Task, submittedData map[string]interface{}) (*Asset, error) {
    asset, err := s.FindAsset(assetId)
    if err != nil {
        return asset, err
    }
    if asset == nil {
        assetError := errors.New("Failed finding an asset with that id.")
        return asset, assetError
    }
    asset.SubmittedData[task.Name] = submittedData
    p := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }

    tasks, _, err := s.FindTasks(p)
    if err != nil {
        return asset, err
    }
    assetVerified := true
    for _, t := range tasks {
        if asset.SubmittedData[t.Name] == nil {
            assetVerified = false
        }
    }
    if assetVerified {
        log.Println("Asset #", asset.Id, "is considered verified!")
    }
    asset.Verified = assetVerified
    _, err = s.EsConn.Index(s.Index, "assets", assetId, nil, asset)
    if err != nil {
        return asset, err
    }
    return asset, nil
}

// CalculateAssetCounts tallies up number of assignments, favorites, etc an asset has and saves it
func (s *Server) CalculateAssetCounts(asset Asset) (Asset, error) {
    assetTmpl := `{
        "query": {
            "bool": {
                "must": [
                {
                    "term": {
                        "assignments.Asset.Id": "%s"
                    }
                }
                ],
                "must_not": [],
                "should": []
            }
        },
        "from": 0,
        "size": 10,
        "sort": [],
        "aggs": {
            "assignments": {
                "terms": {
                    "field": "State",
                    "size": 0
                }
            }
        }
    }`
    assignmentQuery := fmt.Sprintf(assetTmpl, asset.Id)
    assignResults, err := s.EsConn.Search(s.Index, "assignments", nil, assignmentQuery)
    if err != nil {
        return asset, err
    }
    var a assignmentAgg
    err = json.Unmarshal(assignResults.Aggregations, &a)
    if err != nil {
        return asset, err
    }

    if len(asset.Counts) <= 0 {
        asset.Counts = Counts{
            "Favorites":   0,
            "Assignments": 0,
            "finished":    0,
            "skipped":     0,
            "unfinished":  0,
        }
    } else {
        asset.Counts = Counts{
            "Assignments": 0,
            "finished":    0,
            "skipped":     0,
            "unfinished":  0,
        }
    }
    total := 0
    for _, b := range a.Assignments.Buckets {
        asset.Counts[b.Key] = b.Count
                total += b.Count
    }
    asset.Counts["Assignments"] = total

    _, err = s.EsConn.Index(s.Index, "assets", asset.Id, nil, asset)
    if err != nil {
        return asset, err
    }
    return asset, nil
}

//
//    Users
//

// FindUsers returns an array of users in the current project, along with pagination meta information
// 'from' and 'size' parameters determine the offset and limit passed to the database.
func (s *Server) FindUsers(p Params) (users []User, m meta, err error) {
    query := elastigo.Search(s.Index).Type("users").Filter(
        elastigo.Filter().Terms("Project", s.ActiveProjectId),
    ).From(p.From).Size(p.Size)
    if p.SortDir == "desc" {
        query = query.Sort(elastigo.Sort(p.SortBy).Desc())
    } else {
        query = query.Sort(elastigo.Sort(p.SortBy).Asc())
    }

    results, err := query.Result(&s.EsConn)

    if err != nil {
        users = make([]User, 0)
        return users, m, nil
    }

    resultCount := results.Hits.Total

    m.Total = resultCount
    m.From, _ = strconv.Atoi(p.From)
    m.Size, _ = strconv.Atoi(p.Size)

    taskParams := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }

    tasks, m, err := s.FindTasks(taskParams)
    for _, hit := range results.Hits.Hits {
        var user User
        rawMessage := hit.Source
        err = json.Unmarshal(*rawMessage, &user)

        if err != nil {
            err = nil
        }
        if len(tasks) > 0 {
            for _, task := range tasks {
                _, ok := user.Counts[task.Id]
                if !ok {
                    user.Counts[task.Id] = 0
                }
            }
        }
        users = append(users, user)
    }
    return
}

// FindUser looks up a user by id. If a matching user isn't found, it will create a new user and return it.
// TODO: make the CreateUser part optional/conditional?
func (s *Server) FindUser(id string) (user *User, err error) {
    if id == "" {
        userData := strings.NewReader(fmt.Sprintf(`{"Project": "%s"}`, s.ActiveProjectId))
        user, err = s.CreateUser(userData)
        if err != nil {
            return nil, err
        }
        return user, nil
    }

    err = s.EsConn.GetSource(s.Index, "users", id, nil, &user)

    if err != nil {
/* -A.G. (Apparent result of hasty patching; must always return either valid User or error)
        var args map[string]interface{}
        userExists, _ := s.EsConn.ExistsBool(s.Index, "users", id, args)
        if !userExists {
            return nil, nil
        }
*/
        return nil, err
    }

    p := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }

    tasks, _, err := s.FindTasks(p)
    if err == nil {
        for _, task := range tasks {
            _, ok := user.Counts[task.Id]
            if !ok {
                user.Counts[task.Id] = 0
            }
        }
    }
    return user, nil
}

// Creates a user based on the JSON body of the request.
func (s *Server) CreateUser(requestBody io.Reader) (user *User, err error) {

    body, err := ioutil.ReadAll(requestBody)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(body, &user)
    if err != nil {
        return nil, err
    }

    user.Project = s.ActiveProjectId
    user.Favorites = userFavorites{}

    user.Counts = Counts{
        "Favorites":      0,
        "Assignments":    0,
        "VerifiedAssets": 0,
    }

    taskParams := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }

    tasks, _, err := s.FindTasks(taskParams)
    if err != nil {
        for _, task := range tasks {
            user.Counts[task.Id] = 0
        }
    }

    // store user in elasticsearch
    // if user.Id is blank, es will generate a new one
    // if user.Id is NOT blank, es will store the user with that id
    result, err := s.EsConn.Index(s.Index, "users", user.Id, nil, user)
    if err != nil {
        return user, err
    }

    // if the user didn't have an autogenerated id, store it now
    if len(user.Id) == 0 {
        user.Id = result.Id
        _, err = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
        if err != nil {
            return user, err
        }
    }

    return user, nil
}

// Creates a user account with a given user id, 
//     called when a user has a {project_id}_user_id but no matching record is found.
// in other words, this method is used in edge cases.
func (s *Server) CreateUserFromMissingCookieValue(userId string) (User, error) {
    var err error

    user := User{
        Id:      userId,
        Project: s.ActiveProjectId,
    }
    user.Favorites = userFavorites{}
    user.Counts = Counts{
        "Favorites":      0,
        "Assignments":    0,
        "VerifiedAssets": 0,
    }

    taskParams := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }

    tasks, _, err := s.FindTasks(taskParams)
    if err != nil {
        for _, task := range tasks {
            user.Counts[task.Id] = 0
        }
    }

    // store user in elasticsearch
    // if user.Id is blank, es will generate a new one
    // if user.Id is NOT blank, es will store the user with that id
    result, err := s.EsConn.Index(s.Index, "users", user.Id, nil, user)
    if err != nil {
        return user, err
    }

    // if the user didn't have an autogenerated id, store it now
    if len(user.Id) == 0 {
        user.Id = result.Id
        _, err = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
        if err != nil {
            return user, err
        }
    }

    return user, nil
}

// Creates a user account with a given ExternalId. This method is used to link user accounts from third
// party/external registration systems into hive.
func (s *Server) CreateExternalUser(externalId string) (User, error) {
    var user User
    user.ExternalId = externalId
    user.Project = s.ActiveProjectId
    user.Favorites = userFavorites{}
    user.Counts = Counts{
        "Favorites":      0,
        "Assignments":    0,
        "VerifiedAssets": 0,
    }

    taskParams := Params{
        From:    "0",
        Size:    "10",
        SortBy:  "Name",
        SortDir: "asc",
    }

    tasks, _, err := s.FindTasks(taskParams)
    if err != nil {
        for _, task := range tasks {
            user.Counts[task.Id] = 0
        }
    }

    // store user in elasticsearch
    // if user.Id is blank, es will generate a new one
    // if user.Id is NOT blank, es will store the user with that id
    result, err := s.EsConn.Index(s.Index, "users", user.Id, nil, user)
    if err != nil {
        return user, err
    }

    // if the user didn't have an autogenerated id, store it now
    if len(user.Id) == 0 {
        user.Id = result.Id
        _, err = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
        if err != nil {
            return user, err
        }
    }
    return user, nil
}

//
//    Assignments
//

// FindAssignments returns an array of assignments in the current project, given task and state, 
//     along with pagination meta information. 
// 'from' and 'size' parameters determine the offset and limit passed to the database.
func (s *Server) FindAssignments(p Params) (assignments []Assignment, m meta, err error) {
    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return
    }

    if !strings.HasPrefix(p.Task, s.ActiveProjectId) && p.Task != "" {
        p.Task = s.ActiveProjectId + "-" + p.Task
    }

    musts := []string{}
    musts = append(musts, fmt.Sprintf(` { "query": { "match": { "Project": "%s" } } }`, s.ActiveProjectId))

    if p.Task != "" {
        musts = append(musts, fmt.Sprintf(`{ "query": { "match": { "Task": "%s" } } }`, p.Task))
    }

    if p.State != "" {
        musts = append(musts, fmt.Sprintf(` { "query": { "match": { "State": "%s" } } }`, p.State))
    }

    searchQuery := `{
        "query": {
            "filtered": {
                "filter": {
                    "bool": {
                        "must": [%s]
                    }
                }
            }
        },
        "from": %s,
        "size": %s,
        "sort": [ { "%s": { "order" : "%s" } } ]
    }`

    searchJson := fmt.Sprintf(searchQuery, strings.Join(musts, ", "), p.From, p.Size, p.SortBy, p.SortDir)
    results, err := s.EsConn.Search(s.Index, "assignments", nil, searchJson)
    if err != nil {
        return
    }

    m.Total = results.Hits.Total
    m.From, _ = strconv.Atoi(p.From)
    m.Size, _ = strconv.Atoi(p.Size)

    for _, hit := range results.Hits.Hits {
        var assignment Assignment
        rawMessage := hit.Source
        err = json.Unmarshal(*rawMessage, &assignment)
        if err != nil {
            return
        }
        assignments = append(assignments, assignment)
    }
    if len(assignments) <= 0 {
        assignments = make([]Assignment, 0)
    }
    return
}

// FindAssignment looks up an assignment by id.
func (s *Server) FindAssignment(id string) (assignment *Assignment, err error) {

    err = s.EsConn.GetSource(s.Index, "assignments", id, nil, &assignment)
    if err != nil {
        return nil, err
    }
    return assignment, nil
}

// CreateAssignment is called by the userAssignmentHandler to generate an assignment for the given user and task,
// picking an eligible asset for that task and user.
func (s *Server) CreateAssignment(taskId string, userId string) (assignment *Assignment, err error) {

    user, _ := s.FindUser(userId)
    if user == nil {
        tmpUser, err := s.CreateUserFromMissingCookieValue(userId)
        if err != nil {
            userError := errors.New("Assignments can't be created without a user: failed creating a new anon user")
            return nil, userError
        }
        user = &tmpUser
    }

    task, err := s.FindTask(taskId)
    if err != nil {
        return nil, err
    }

    if task.CurrentState != "available" {
        taskError := errors.New("Invalid task")
        return nil, taskError
    }

    searchQuery := `{
        "query": {
            "bool": {
                "must": [
                {
                    "term": {
                        "assignments.Project": "%s"
                    }
                },
                {
                    "term": {
                        "assignments.Task": "%s"
                    }
                },
                {
                    "term": {
                        "assignments.User": "%s"
                    }
                },
                {
                    "term": {
                        "assignments.State": "unfinished"
                    }
                }
                ]
            }
        }
    }`

    searchJson := fmt.Sprintf(searchQuery, s.ActiveProjectId, taskId, userId)

    results, err := s.EsConn.Search(s.Index, "assignments", nil, searchJson)
    if err != nil {
        return nil, err
    }

    if results.Hits.Total > 0 {
        // found an unfinished assignment
        err = json.Unmarshal(*results.Hits.Hits[0].Source, &assignment)
        if err != nil {
            return nil, err
        }
        return assignment, nil

    } else {
        // create a new assignment
        assignmentAsset, err := s.FindAssignmentAsset(*task, *user)
        if err != nil {
            return nil, err
        }

        // Set counts on asset
        if len(assignmentAsset.Counts) <= 0 {
            assignmentAsset.Counts = Counts{
                "Favorites":   0,
                "Assignments": 0,
                "finished":    0,
                "skipped":     0,
                "unfinished":  0,
            }
        }

        // Since this asset is being assigned now, update the total assignments count
        assignmentAsset.Counts["Assignments"] += 1

        // And update the unfinished count, since it's a new assignment
        assignmentAsset.Counts["unfinished"] += 1

        _, err = s.EsConn.Index(s.Index, "assets", assignmentAsset.Id, nil, assignmentAsset)
        if err != nil {
            return nil, err
        }

        assignmentId := strings.Join([]string{s.ActiveProjectId, taskId, assignmentAsset.Id, user.Id}, "HIVE")
        assignment = &Assignment{
            Id:      assignmentId,
            User:    userId,
            Project: s.ActiveProjectId,
            Task:    taskId,
            Asset:   assignmentAsset,
            State:   "unfinished",
        }

        _, err = s.EsConn.Index(s.Index, "assignments", assignment.Id, nil, assignment)
        if err != nil {
            return nil, err
        }
        return assignment, nil
    }
}

// FindAssignmentAsset returns an eligible asset for a given task and user, basing this on AssignmentCriteria.
// It is called from CreateAssignment.
func (s *Server) FindAssignmentAsset(task Task, user User) (Asset, error) {
    var assignmentAsset Asset
    var assetIds []string

    assetQuery := fmt.Sprintf(`{
        "query": {
            "bool": {
                "must": [
                {
                "term": {
                  "assignments.Task": "%s"
                }
                      },
              {
          "term": {
            "assignments.User": "%s"
          }
                },
                {
                    "term": {
                        "assignments.Project": "%s"
                    }
                }
                ]
            }
        },
        "from": 0,
        "size": %d
    }`, task.Id, user.Id, s.ActiveProjectId, user.Counts["Assignments"])
    assetResults, err := s.EsConn.Search(s.Index, "assignments", nil, assetQuery)
    if err != nil {
        return assignmentAsset, err
    }
    for _, hit := range assetResults.Hits.Hits {
        idParts := strings.Split(hit.Id, "HIVE")
        assetIds = append(assetIds, idParts[2])
    }

    // the parts of a 'bool' query - so far no need for 'should'
    musts := []string{}
    mustNots := []string{}

    // build up the pieces of the full elasticsearch query
    for taskName, ruleI := range task.AssignmentCriteria.SubmittedData {
        rule := ruleI.(map[string]interface{})

        if len(rule) == 0 {
            // an empty rule means assets should have no data submitted for this task
            tmpl := `{
                "missing": {
                    "field": "SubmittedData.%s"
                }
            }`

            musts = append(musts, fmt.Sprintf(tmpl, task.Name))

        } else {
            // assets must have data submitted that exactly matches the rule
            for fieldName, fieldValue := range rule {
                tmpl := `{
                    "query": {
                        "match": {
                            "SubmittedData.%s.%s": "%s"
                        }
                    }
                }`
                musts = append(musts, fmt.Sprintf(tmpl, taskName, fieldName, fieldValue))
            }
        }
    }

    // limit query results to assets in this project
    projectTmpl := `{
        "query": {
            "match": {
                "Project": "%s"
            }
        }
    }`
    musts = append(musts, fmt.Sprintf(projectTmpl, s.ActiveProjectId))

    if len(assetIds) > 0 {
        assetTmpl := `{ "query": { "terms": { "Id": [ %s ] } } }`
        assetIdString := "\"" + strings.Join(assetIds, "\",\"") + "\""
        mustNots = append(mustNots, fmt.Sprintf(assetTmpl, assetIdString))
    }

    mustsJson := strings.Join(musts, ", ")
    mustNotsJson := strings.Join(mustNots, ", ")

    var args map[string]interface{}
    matchAllQuery := `{ "query": { "match_all" : { } } }`
    countResponse, err := s.EsConn.Count(s.Index, "assets", args, matchAllQuery)
    if err != nil {
        return assignmentAsset, err
    }

    // finally, compose the entire filtered query
    searchQuery := fmt.Sprintf(
        `{"query": {
            "filtered": {
                "filter": {
                    "bool": {
                        "must": [%s],
                        "must_not":[%s]
                    }
                }
            }
        },
        "from": 0,
        "size": %d
    }`, mustsJson, mustNotsJson, countResponse.Count)

    results, err := s.EsConn.Search(s.Index, "assets", nil, searchQuery)
    if err != nil {
        return assignmentAsset, err
    }

    if results.Hits.Total <= 0 {
        err = errors.New("No assets found")
        return assignmentAsset, err

    } else {
        randomHit := rand.Intn(len(results.Hits.Hits))
        rawMessage := results.Hits.Hits[randomHit].Source
        err = json.Unmarshal(*rawMessage, &assignmentAsset)
        if err != nil {
            return assignmentAsset, err
        }
    }
    return assignmentAsset, nil
}

// CreateAssetAssignment is called by the AssignAssetHandler to generate
//     a new assignment for a particular asset, task and user
func (s *Server) CreateAssetAssignment(
        taskId string, userId string, assetId string) (assignment *Assignment, err error) {
    user, _ := s.FindUser(userId)
    if user == nil {
        tmpUser, err := s.CreateUserFromMissingCookieValue(userId)
        if err != nil {
            userError := errors.New(
                "Assignments can't be created without a user: failed creating a new anon user")
            return nil, userError
        }
        user = &tmpUser
    }

    asset, err := s.FindAsset(assetId)
    if asset == nil {
        assetError := errors.New("Failed finding an asset with that id.")
        return nil, assetError
    }

    // Set counts on asset
    if len(asset.Counts) <= 0 {
        asset.Counts = Counts{
            "Favorites":   0,
            "Assignments": 0,
            "finished":    0,
            "skipped":     0,
            "unfinished":  0,
        }
    }
    asset.Counts["Assignments"] += 1
    asset.Counts["unfinished"] += 1
    _, err = s.EsConn.Index(s.Index, "assets", asset.Id, nil, asset)
    if err != nil {
        log.Println(err)
    }

    assignmentId := strings.Join([]string{s.ActiveProjectId, taskId, assetId, userId}, "HIVE")
    assignment = &Assignment{
        Id:      assignmentId,
        User:    userId,
        Project: s.ActiveProjectId,
        Task:    taskId,
        Asset:   *asset,
        State:   "unfinished",
    }

    _, err = s.EsConn.Index(s.Index, "assignments", assignment.Id, nil, assignment)
    if err != nil {
        return nil, err
    }
    return assignment, nil
}

func (s *Server) UpdateAssignment(requestBody io.Reader) (assignment *Assignment, err error) {
    body, err := ioutil.ReadAll(requestBody)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(body, &assignment)
    if err != nil {
        return nil, err
    }

    asset, _ := s.FindAsset(assignment.Asset.Id)
    if asset != nil {
        // Set counts on asset
        if len(asset.Counts) <= 0 {
            asset.Counts = Counts{
                "Favorites":   0,
                "Assignments": 1,
                "finished":    0,
                "skipped":     0,
                "unfinished":  1,
            }
        }

        asset.Counts[assignment.State] += 1
        asset.Counts["unfinished"] -= 1

        _, err = s.EsConn.Index(s.Index, "assets", asset.Id, nil, asset)
        if err != nil {
            return nil, err
        }
        // ensure the asset is updated on the assignment record
        assignment.Asset = *asset
    }

    _, err = s.EsConn.Index(s.Index, "assignments", assignment.Id, nil, assignment)
    if err != nil {
        return nil, err
    }
    // refresh the index, attempting to fix "skipped" assignment issue #4
    _, err = s.EsConn.Refresh(s.Index)
    if err != nil {
        return nil, err
    }

    // add finished assignments to the user's list
    if assignment.State == "finished" {
        user, err := s.FindUser(assignment.User)
        if err != nil {
            return nil, err
        }
        user.Counts["Assignments"]++
        user.Counts[assignment.Task]++

        p := Params{
            From:    "0",
            Size:    "10",
            SortBy:  "Name",
            SortDir: "asc",
        }

        tasks, _, err := s.FindTasks(p)
        if err != nil {
            for _, task := range tasks {
                // Set any missing task counts to zero
                _, ok := user.Counts[task.Id]
                if !ok {
                    user.Counts[task.Id] = 0
                }
            }
        }

        _, err = s.EsConn.Index(s.Index, "users", user.Id, nil, user)
        if err != nil {
            return nil, err
        }
    }
    return assignment, nil
}

// CountAssignments returns a map of assignment states to totals for each scoped to the current project.
func (s *Server) CountAssignments() (assignmentCount map[string]int, err error) {
    projectQuery := fmt.Sprintf(`{
        "aggs": {
            "assignments": {
                "terms": {
                    "field": "State",
                    "size": 0
                }
            }
        },
        "query": {
            "filtered": {
                "filter": {
                    "bool": {
                        "must": [
                        {
                            "query": {
                                "match": {
                                    "Project": "%s"
                                }
                            }
                        }
                        ]
                    }
                }
            }
        }
    }`, s.ActiveProjectId)
    results, err := s.EsConn.Search(s.Index, "assignments", nil, projectQuery)
    if err != nil {
        return
    }
    var a assignmentAgg
    err = json.Unmarshal(results.Aggregations, &a)
    if err != nil {
        return nil, err
    }

    assignmentCount = make(map[string]int)
        total := 0
    for _, b := range a.Assignments.Buckets {
        assignmentCount[strings.Title(b.Key)] = b.Count
                total += b.Count
    }
    assignmentCount["Total"] = total
    return assignmentCount, nil
}

