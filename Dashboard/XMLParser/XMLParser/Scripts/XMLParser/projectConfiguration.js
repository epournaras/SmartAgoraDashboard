$(document).ready(function() {
    $('.validating-fields').on('input', function() {
        var input = $(this);
        var is_name = input.val();
        
        if(is_name){
            input.removeClass("parsley-error");
        }
        else{
            input.addClass("parsley-error");
        }
    });
 });

$(document).ready(function () {

    loadProjects();
    //--------------------------------------
    //Scripting about project creation form
    //--------------------------------------
    $("#projectResultDiv").hide();
    $("#btnCreateProject").on('click', function () {
        var project_txtProjectName = $("#project_txtProjectName").val();
        var project_txtProjectDesc = $("#project_txtProjectDesc").val();
        var autoAssingmentSts = $("#autoAssignment").is(":checked");
        var helpTextStr = $("#helpText").val();

        project_txtProjectDesc = project_txtProjectDesc + "#" + helpTextStr + "#" + autoAssingmentSts;
        var project_txtProjectId = "";
        var project_alertMsgText = "please add following in Project creation section :</br>";
        if (project_txtProjectName === "") {
            $("#project_txtProjectName").addClass("parsley-error");
            project_alertMsgText += "=> Project Name</br>";
        }


        if (project_alertMsgText != "please add following in Project creation section :</br>") {
            project_alertMsgText = "Before creating project, " + project_alertMsgText;
            new PNotify({
                title: 'Info!',
                text: project_alertMsgText,
                type: 'info',
                delay: 3500
            });
            return;
        }
        if (project_alertMsgText === "please add following in Project creation section :</br>") {
            $.ajax({
                type: "GET",
                url: createProjectUrl,
                data: { projectId: project_txtProjectId, projectName: project_txtProjectName, projectDesc: project_txtProjectDesc},
                contentType: "application/json",
                dataType: "json",
                beforeSend: function () {
                    $("#overlay").css("display", "block");
                    NProgress.start();
                },
                complete: function () {
                    NProgress.done();
                    $("#overlay").css("display", "none");
                },
                success: function (data) {
                    if (data === "LogOut") {
                        new PNotify({
                            title: 'Info!',
                            text: 'Please Login to create Project!',
                            type: 'info',
                            delay: 3500,
                            after_close: function () {
                                window.location.href = '/Home/index';
                            }
                        });
                    }
                    else {
                        new PNotify({
                            title: 'Success!',
                            text: 'Project Created Successfully',
                            type: 'success',
                            delay: 3500,
                        });
                        loadProjects();
                        clearProjectCreationForm();
                        var json = JSON.stringify(data.aaData, undefined, 2);
                        document.getElementById('result').innerHTML = json;
                        //$("#projectResultDiv").show();
                    }
                },
                error: function (data, textStatus, jqXHR) {
                    new PNotify({
                        title: 'Error!',
                        text: 'Some error occured, please try again',
                        type: 'error',
                        delay: 3000
                    });
                    clearProjectCreationForm();
                }
            });
        }
    });
    
     //--------------------------------------
    //Scripting about project deletion form
    //--------------------------------------
 
    $("#btnDeleteProject").on('click', function () {
        var project_txtProjectIdVal = $("#project_txtProjectId").val();
     
        var project_alertMsgText = "please add following in Project creation section :</br>";
        if (project_txtProjectIdVal === "") {
            $("#project_txtProjectId").addClass("parsley-error");
            project_alertMsgText += "=> Project Name</br>";
        }

        if (project_alertMsgText === "please add following in Project creation section :</br>") {
            $.ajax({
                type: "GET",
                url: deleteProjectUrl,
                data: { projectId: project_txtProjectIdVal},
                contentType: "application/json",
                dataType: "json",
                beforeSend: function () {
                    $("#overlay").css("display", "block");
                    NProgress.start();
                },
                complete: function () {
                    NProgress.done();
                    $("#overlay").css("display", "none");
                },
                success: function (data) {
                    if (data === "LogOut") {
                        new PNotify({
                            title: 'Info!',
                            text: 'Please Login to delete Project!',
                            type: 'info',
                            delay: 3500,
                            after_close: function () {
                                window.location.href = '/Home/index';
                            }
                        });
                    }
                    else {
                        new PNotify({
                            title: 'Success!',
                            text: 'Project Deleted Successfully',
                            type: 'success',
                            delay: 3500,
                        });
                       
                        clearProjectDeletionForm();
                        
                        //$("#projectResultDiv").show();
                    }
                },
                error: function (data, textStatus, jqXHR) {
                    new PNotify({
                         title: 'Success!',
                            text: 'Project Deleted Successfully',
                            type: 'success',
                            delay: 3500,
                    });
                    clearProjectDeletionForm();
                }
            });
        }
    });


    //-----------------------------------
    //Scripting about task creation form
    //-----------------------------------
    $('#task_drpProject').multiselect({
        enableFiltering: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true,
        nonSelectedText: 'Select Project',
    });

    $("#taskResultDiv").hide();

    //to preventing from space in input that cuases issues on android side.
    $("#task_txtTaskName").keypress(function (e) {
        if (e.which === 32)
            return false;
    });
    //$("#task_txtTotal, #task_txtMatching").on('keyup', function () {
    //    debugger;
    //    if (((event.keyCode !== 13 || event.which !== 13)
    //        && (event.keyCode !== 8 || event.which !== 8)
    //        && (event.keyCode !== 9 || event.which !== 9))
    //        && (isNaN(parseInt($(this).val())))) {
    //        new PNotify({
    //            title: 'Info!',
    //            text: "Please enter a valid value",
    //            type: 'info',
    //            delay: 2000,
    //            nonblock: {
    //                nonblock: true
    //            }
    //        });
    //        var id = $(this).prop("id");
           
    //        $("#" + id).val("");
    //    }
    //});
    $("#btnCreateTask").on("click", function createTask() {
        
        var task_projectId = $("#task_drpProject").val();
        var task_name = $("#task_txtTaskName").val();
        var task_desc = $("#task_txtTaskDesc").val();
        var task_state = $("#task_drpState").val();
        //var task_total = $("#task_txtTotal").val();
        //var task_matching = $("#task_txtMatching").val();

        var task_alertMsg = "Please add following in Task Creation Section:</br>";
        if (task_projectId == "" || task_projectId == "0") {
            task_alertMsg += "=> Project</br>";
             $("#task_drpProject").addClass("parsley-error");
        }
        if (task_name == "") {
            task_alertMsg += "=> Task Name</br>";
             $("#task_txtTaskName").addClass("parsley-error");
        }
        if (task_desc == "") {
            task_alertMsg += "=> Task Description</br>";
            $("#task_txtTaskDesc").addClass("parsley-error");
        }
        if (task_state == "") {
            task_alertMsg += "=> Task State</br>";
        }
        //if (task_total == "") {
        //    task_alertMsg += "=> Completion Criteria Total</br>";
        //    $("#task_txtTotal").addClass("parsley-error");
        //}
        //if (task_matching == "") {
        //    task_alertMsg += "=> Completion Criteria Matching</br>";
        //    $("#task_txtMatching").addClass("parsley-error");
        //}
        if (task_alertMsg != "Please add following in Task Creation Section:</br>") {
            task_alertMsg = "Before creating task, " + task_alertMsg;
            new PNotify({
                title: 'Info!',
                text: task_alertMsg,
                type: 'info',
                delay: 3500
            });
            return;
        }
        if (task_alertMsg === "Please add following in Task Creation Section:</br>") {
            $.ajax({
                type: "GET",
                url: createTaskUrl,
                data: { projectId: task_projectId, name: task_name, desc: task_desc, state: task_state/*, total: task_total, matching: task_matching*/ },
                contentType: "application/json",
                dataType: "json",
                beforeSend: function () {
                    $("#overlay").css("display", "block");
                    NProgress.start();
                },
                complete: function () {
                    NProgress.done();
                    $("#overlay").css("display", "none");
                },
                success: function (data) {
                    if (data === "LogOut") {
                        new PNotify({
                            title: 'Info!',
                            text: "Please Login to create Task! ",
                            type: 'info',
                            delay: 3500,
                            after_close: function () {
                                var url = "/Home/Index";
                                window.location.href = url;
                            }
                        });
                    }

                    else {
                        new PNotify({
                            title: 'Success',
                            text: 'Task Created Successfully!',
                            type: 'success',
                            delay: 3500
                        });
                        clearTaskCreationForm();
                        var json = JSON.stringify(data.aaData, undefined, 2);
                        document.getElementById('result').innerHTML = json;
                        //$("#taskResultDiv").show();
                    }
                },
                error: function (data, textStatus, jqXHR) {
                    new PNotify({
                        title: 'Error!',
                        text: "Some error occured, please try again",
                        type: 'error',
                        delay: 3500
                    });
                    clearTaskCreationForm();
                }
            });
        }

    });
    


    //-----------------------------------
    //Scripting about Assignment creation form
    //-----------------------------------
    assignment_assetmsg = "";
    assignment_usermsg = "";
    assignment_taskmsg = "";
    assignment_projectId = "";
    assignment_taskId = "";
    assignment_assetId = "";
    assignment_stsTasks = false;
    assignment_stsAssets = false;
    assignment_stsUser = false;
    $("#assignmentResultDiv").hide();
    $('#assignment_drpUser').multiselect({
        enableFiltering: true,
        includeSelectAllOption: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true,
        selectAllValue: 'select-all-value',
    });

    $('#assignment_drpAsset').multiselect({
        enableFiltering: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true,
        nonSelectedText: 'Select Asset',
        onChange: function (element, checked) {
            $('#assignment_drpAsset').multiselect('updateButtonText', true);
            var assignment_assets = $('#assignment_drpAsset option:selected');
            if (assignment_assets != null) {
                if (assignment_assets[0].value != 0)
                    assignment_assetId = assignment_assets[0].value;
                else {
                   assignment_assetId = "";
                }
            }
            else
                assignment_resetAssets();
        }
    });

    $('#assignment_drpTask').multiselect({
        enableFiltering: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true,
        nonSelectedText: 'Select Task',
        onChange: function (element, checked) {
            $('#assignment_drpTask').multiselect('updateButtonText', true);
            var assignment_tasks = $('#assignment_drpTask option:selected');
            if (assignment_tasks != null) {
                if (assignment_tasks[0].value != 0)
                    assignment_taskId = assignment_tasks[0].value;
                else {
                    assignment_taskId = "";
                }
            }
            else
                assignment_resetTasks();
        }
    });

    $('#assignment_drpProject').multiselect({
        enableFiltering: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true,
        nonSelectedText: 'Select Project',
        onChange: function (element, checked) {
            
            $('#assignment_drpProject').multiselect('updateButtonText', true);
            assignment_projects = $('#assignment_drpProject option:selected');
            if (assignment_projects != null) {
                if (assignment_projects[0].value != 0) {
                    assignment_projectId = assignment_projects[0].value;
                    $("#overlay").css("display", "block");
                    NProgress.start();
                    assignment_loadAssets(assignment_projectId);
                    assignment_loadUsers(assignment_projectId);
                    assignment_loadTasks(assignment_projectId);
                }
                else {
                    assignment_resetAssets();
                    assignment_resetTasks();
                    assignment_resetUsers();
                }
            }
            else {
                assignment_resetAssets();
                assignment_resetTasks();
                assignment_resetUsers();
            }
        }
    });

    $('#assignment_drpUser').multiselect('disable');
    $('#assignment_drpTask').multiselect('disable');
    $('#assignment_drpAsset').multiselect('disable');
                  
    
    $("#btnCreateAssignment").on('click',function(){

        var assignment_users = $('#assignment_drpUser option:selected');
        var assignment_userId = [];
        $(assignment_users).each(function (index, brand) {
            assignment_userId.push([$(this).val()]);
        });

        var assignment_alertMsg = "please select following in Assignment creation section :</br>";
        if (assignment_projectId == "" || assignment_projectId == null || assignment_projectId=="0") {
            assignment_alertMsg += "=> Project</br>";
        } 
    if (assignment_userId == "" || assignment_userId == null || assignment_userId=="0") {
            assignment_alertMsg += "=> User(s)</br>";
        }   
    if (assignment_assetId == "" || assignment_assetId == null || assignment_assetId=="0") {
            assignment_alertMsg += "=> Asset</br>";
        }   
    if (assignment_taskId == "" || assignment_taskId == null || assignment_taskId=="0") {
            assignment_alertMsg += "=> Task</br>";
        }   
        if (assignment_alertMsg != "please select following in Assignment creation section :</br>") {
            assignment_alertMsg = "Before creating assignment, " + assignment_alertMsg;
            new PNotify({
                title: 'Info!',
                text: assignment_alertMsg,
                type: 'info',
                delay: 3500
            })
            return;
        }

        $.ajax({
            type: "GET",
            url: createAssignmentUrl,
            data: { projectId: assignment_projectId, taskId: assignment_taskId, assetId: assignment_assetId, userIdsList: JSON.stringify(assignment_userId) },
            contentType: "application/json",
            dataType: "json",
            beforeSend: function () {
                $("#overlay").css("display", "block");
                NProgress.start();
            },
            complete: function () {
                NProgress.done();
                $("#overlay").css("display", "none");
            },
            success: function (data) {
                if (data === "LogOut") {
                    new PNotify({
                        title: 'Info!',
                        text: "Please Login to create Assignment! ",
                        type: 'info',
                        delay: 3500,
                        after_close: function () {
                            window.location.href = '/Home/index';
                        }
                    });
                }
                else {
                    new PNotify({
                        type: 'success',
                        title: 'Assignment Created Successfully!',
                        type: 'success',
                        delay: 3500
                    });
                    resetAll();
                    //var json = JSON.stringify(data.aaData, undefined, 2);
                    //document.getElementById('result').innerHTML = json;
                    //$("#assignmentResultDiv").show();
                } 
            },
            error: function (data, textStatus, jqXHR) {
                new PNotify({
                    //title: 'Error!',
                    title: '',
                    text: "Assignment already assignned!",
                    //text: "Some error occured, please try again",
                    type: 'error',
                    delay: 3500
                });
            }
        });
    });
});


function assignment_resetProjects() {
    assignment_projectId = "";
    $('#assignment_drpProject option:selected').prop('selected', false);
    $('#assignment_drpProject').multiselect('refresh');
}
function assignment_resetTasks() {
    assignment_taskId = "";
    $('#assignment_drpTask option:selected').prop('selected', false);
    $('#assignment_drpTask').multiselect('refresh');
    $('#assignment_drpTask').multiselect('disable');
}
function assignment_resetUsers() {
    $('#assignment_drpUser option:selected').prop('selected', false);
    $('#assignment_drpUser').multiselect('refresh');
    $('#assignment_drpUser').multiselect('disable');

}
function assignment_resetAssets() {
    assignment_assetId = "";
    $('#assignment_drpAsset option:selected').prop('selected', false);
    $('#assignment_drpAsset').multiselect('refresh');
    $('#assignment_drpAsset').multiselect('disable');
}

function resetAll() {
    assignment_resetProjects();
    assignment_resetAssets();
    assignment_resetTasks();
    assignment_resetUsers();
}
    

function clearTaskCreationForm() {
    $('#task_drpProject option:selected').prop('selected', false);
    $('#task_drpProject').multiselect('refresh');
    $("#task_txtTaskName").val("");
    $("#task_txtTaskDesc").val("");
    $("#task_drpState").attr('selectedIndex', 0);
    $("#task_txtTotal").val("");
    $("#task_txtMatching").val("");


}

function clearProjectCreationForm() {
    $("#project_txtProjectName").val("");
    $("#project_txtProjectDesc").val("");
    $("#helpText").val("");
}
function clearProjectDeletionForm() {
    $("#project_txtProjectId").val("");
    
}


function loadProjects() {
    
    $.ajax({
        type: "POST",
        url: loadProjectIdsUrl,
        contentType: "application/json",
        dataType: "json",
        beforeSend: function () {
            $("#overlay").css("display", "block");
            NProgress.start();
        },
        complete: function () {
            NProgress.done();
            $("#overlay").css("display", "none");
        },
        success: function (data) {
            
            var options = [];
            options.push({ label: "Select Project", title: "Select Project", value: 0 });
            for (var i = 0; i < data.aaData.length; i++) {
                options.push({ label: data.aaData[i].Name, title: data.aaData[i].Name, value: data.aaData[i].Id });
            }

            $('#assignment_drpProject').multiselect('dataprovider', options);
            $('#assignment_drpProject').multiselect('enable');

            $('#task_drpProject').multiselect('dataprovider', options);
            $('#task_drpProject').multiselect('enable');


        },
        error: function (data, textStatus, jqXHR) {
            new PNotify({
                title: 'Error!',
                text: "Load Project.Some error occured, please try again",
                type: 'error',
                delay: 3500
            });
        }
    });
}


function assignment_loadUsers(projectId) {
    $.ajax({
        type: "GET",
        url: loadUsersUrl,
        data: { projectId: projectId },
        contentType: "application/json",
        dataType: "json",
        complete: function () {

            if (assignment_stsAssets && assignment_stsTasks && assignment_stsUser) {
                $("#overlay").css("display", "none");
                NProgress.done();
                if (assignment_assetmsg != "" || assignment_taskmsg != "" || assignment_usermsg != "") {
                    new PNotify({
                        title: 'Info!',
                        text: assignment_taskmsg + assignment_assetmsg + assignment_usermsg,
                        type: 'info',
                        delay: 3500
                    })
                }
                assignment_stsAssets = false;
                assignment_assetmsg = "";
                assignment_stsTasks = false;
                assignment_taskmsg = "";
                assignment_stsUser = false;
                assignment_usermsg = ""

            }
        },
        success: function (data) {
            
            if (data.userFoundMessage === "nouserfound") {
                assignment_resetUsers();
                assignment_usermsg = "No User found in this Project</br>";
            }
            else {
                var options = [];
                //options.push({ label: "Select User", title: "Select User", value: 0 });
                for (var i = 0; i < data.aaData.length; i++) {
                    options.push({ label: data.aaData[i].Name, title: data.aaData[i].Name, value: data.aaData[i].Id });
                }
                $('#assignment_drpUser').multiselect('dataprovider', options);
                $('#assignment_drpUser').multiselect('enable');
            }
            assignment_stsUser = true;
        },
        error: function (data, textStatus, jqXHR) {
            new PNotify({
                title: 'Error!',
                text: "Load Users.Some error occured, please try again",
                type: 'error',
                delay: 3500
            });
        }
    });
}


function assignment_loadAssets(projectId) {
    $.ajax({
        type: "GET",
        url: loadAssetsUrls,
        data: { projectId: projectId },
        contentType: "application/json",
        dataType: "json",
        complete: function () {
            if (assignment_stsAssets && assignment_stsTasks && assignment_stsUser) {
                $("#overlay").css("display", "none");
                NProgress.done();
                if (assignment_assetmsg != "" || assignment_taskmsg != "" || assignment_usermsg != "") {
                    new PNotify({
                        title: 'Info!',
                        text: assignment_taskmsg + assignment_assetmsg + assignment_usermsg,
                        type: 'info',
                        delay: 3500
                    })
                }
                assignment_stsAssets = false;
                assignment_assetmsg = ""
                
                assignment_stsTasks = false;
                assignment_taskmsg = "";
                
                assignment_stsUser = false;
                assignment_usermsg = ""
            }
        },
        success: function (data) {
            
            if (data.assetFoundMessage === "noAssetfound") {
                assignment_resetAssets();
                assignment_assetmsg = "No Asset found in this Project</br>";
            }
            else {
                $("#assignment_drpAsset").empty();

                var options = [];
                options.push({ label: "Select Asset", title: "Select Asset", value: 0 });
                for (var i = 0; i < data.aaData.length; i++) {
                    options.push({ label: data.aaData[i].Name, title: data.aaData[i].Name, value: data.aaData[i].Id });
                }

                $('#assignment_drpAsset').multiselect('dataprovider', options);
                $('#assignment_drpAsset').multiselect('enable');
            }
            assignment_stsAssets = true;
        },
        error: function (data, textStatus, jqXHR) {
            new PNotify({
                title: 'Error!',
                text: "Load Asset.Some error occured, please try again",
                type: 'error',
                delay: 3500
            });
        }
    });
}

function assignment_loadTasks(projectId) {
    $.ajax({
        type: "GET",
        url: loadTaskIdsUrl,
        data: { projectId: projectId },
        contentType: "application/json",
        dataType: "json",
        complete: function () {
            
            if (assignment_stsAssets && assignment_stsTasks && assignment_stsUser === true) {
                $("#overlay").css("display", "none");
                NProgress.done();
                if (assignment_assetmsg != "" || assignment_taskmsg != "" || assignment_usermsg != "") {
                    new PNotify({
                        title: 'Info!',
                        text: assignment_taskmsg + assignment_assetmsg + assignment_usermsg,
                        type: 'info',
                        delay: 3500
                    })

                }
                assignment_stsAssets = false;
                assignment_assetmsg = ""

                assignment_stsTasks = false;
                assignment_taskmsg = "";
                
                assignment_stsUser = false;
                assignment_usermsg = ""
            }
        },
        success: function (data) {
            
            if (data.taskFoundMessage === "notaskfound") {
                assignment_resetTasks();
                assignment_taskmsg = "No Task found in this Project</br>";

            }
            else {
                $("#assignment_drpTask").empty();

                var options = [];
                options.push({ label: "Select Task", title: "Select Task", value: 0 });
                for (var i = 0; i < data.aaData.length; i++) {
                    options.push({ label: data.aaData[i].Name, title: data.aaData[i].Name, value: data.aaData[i].Id });
                }

                $('#assignment_drpTask').multiselect('dataprovider', options);
                $('#assignment_drpTask').multiselect('enable');
            }
            assignment_stsTasks = true;
        },
        error: function (data, textStatus, jqXHR) {
            new PNotify({
                title: 'Error!',
                text: "Load Task.Some error occured, please try again",
                type: 'error',
                delay: 3500
            });
        }
    });
}