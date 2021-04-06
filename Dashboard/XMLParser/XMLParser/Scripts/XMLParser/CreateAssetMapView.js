var startAndDestinationModal = {}; //start and destination modal
var mainModel = {}; // complete asset modal
var objMainModel = {}; //Containing "Start and Destination Modal" and "SampleDataModel"
var lstQuestionDataModel = []; // List of Questions

var objSensorId = 0; //sensor id
var lstSensors = []; // List of sensor related to a question

var objQuestionId = 0; //question id
var addedQuestionCount = 0; //number of question that has been added

var objOptionId = 1; //option id
var objOption = {}; // option object
var lstOptions = []; //list of option related to a question

var objCombinationId = 1; //combination id
var objCombination = {}; //combination object
var lstCombinations = []; // list of combination related to a question

var assetMode = null; //Selected mode of asset
var questionType = null;

var optionsAdded = 0; //number of options that has been added
var optionsString = null;
var optionsCombinations = [];
var selectOptionString = null;
var combinationFound = false;
var sequence = 1;
var removeCombination = false;
var lstOptionlen = 0;
var questionLocationMarkers = [];
var boolAssociateQuestion = false;
var projectIds = [];
var updateSts = false;
var editedQuestionId = 0;
var questionAddressList = [];
var optionCounter = 0;
var combinationCounter = 0;
var questionOldLat = "";
var questionOldLng = "";

// makes combination string for example is str
// is "12" then combination will be "1,2,12" 
var combinations = function (str) {
    var lenStr = str.length;
    var result = [];
    var indexCurrent = 0;
    while (indexCurrent < lenStr) {
        var char = str.charAt(indexCurrent);
        var x;
        var arrTemp = [char];
        for (x in result) {
            arrTemp.push("" + result[x] + char);
        }
        result = result.concat(arrTemp);
        indexCurrent++;
    }
    return result;
}

// initialize dropdowns(project, question type, sensor, frequency),
// calls loadProject, loadSensors and loadQuestion function, 
// number validation for defaultCredit, time and vicinity and
// set Focused address input(question address, start and destination Address and grouppoint)
$(document).ready(function () {
    //**************************
    //Multiselect initialization
    //**************************
    $('#drpProject').multiselect({
        enableFiltering: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true,
        nonSelectedText: 'Select Project',
    });

    $('#questionType').multiselect({
        enableFiltering: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true
    });

    $("#dropdownSensors").multiselect({
        enableFiltering: true,
        includeSelectAllOption: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true,
        selectAllValue: 'select-all-value',
        nonSelectedText: 'Select Sensor(s)'
    });

    $('#txtFrequency').multiselect({
        enableFiltering: true,
        buttonWidth: '257px',
        buttonClass: 'form-control',
        inheritClass: true
    });



    sequence = 1;
    $("#txtSequence").val(sequence);
    LoadProjects();
    LoadQuestions(lstQuestionDataModel);
    LoadSensors();

    //************************************************************************
    // checking default credit, time and vicinity
    // to verify the input must be positive integer
    //and disabling next button if default credit is empty
    //************************************************************************
    $(".chk-number-input").on('keyup focusout', function () {
        if ((event.keyCode !== 13 || event.which !== 13) && (isNaN(parseInt($(this).val())) || parseInt($(this).val()) == 0)) {
            new PNotify({
                title: 'Info!',
                text: "Please enter a valid value",
                type: 'info',
                delay: 2000,
                nonblock: {
                    nonblock: true
                }
            });
            var id = $(this).prop("id");
            if (id === "defaultCredit")
                $(".buttonNext").addClass("buttonDisabled");
            $("#" + id).val("");
        }
        else {
            $(".buttonNext").removeClass("buttonDisabled");

        }
    });

    //************************************************************************
    // Getting focused Address input
    //************************************************************************
    $("#startAddress, #destinationAddress, #questionAddress, #DIASGrouppointAddress").on('click', function () {
        focusedInput = $(this).prop('id');
    });

});

//************************************************************************
// saving current view's values before going to next on next button click
// and set next step
//************************************************************************


function saveProjectIdAndAssetMode() {
    optionCounter = 0;
    combinationCounter = 0;
    mainModel.projectId = $("#drpProject").val();
    startAndDestinationModal.Mode = $("#assetMode").val();

    $("#projectModeModal").modal('hide');
    if (startAndDestinationModal.Mode === "Decision" || startAndDestinationModal.Mode === "Sequence") {
        setStartAndDestinationStep();
        $("#questionModelBtn").hide();
        $("#startDestinaitonModalBtn").show();
        $("#startDestinationModal").modal('show');
    }
    else {
        setQuestionStep();
        $("#startDestinaitonModalBtn").hide();
        $("#questionModelBtn").show();
        $("#questionDetailModal").modal('show');
    }
}

function setStartAndDestinationStep() {
    focusedInput = null;
    $("#startAddressDiv").show();
    $("#destinationAddressDiv").show();
    $("#googleMap").removeClass("disabledMap");
    focusedInput = "startAddress";
    $("#checkStartDest").show();
    $("#startAddressRadioBtn").iCheck('check');

    $("#DIASGrouppointAddressRadioBtnOnMap").hide();
    $("#DIASGrouppointAddressRadioBtn").iCheck('uncheck');

    $("#QuestionRadioBtnOnMap").hide();
    $("#questionAddressRadioBtn").iCheck('check');
}

function saveStartAndDestinationData() {
    var msg = "";
    
    var assetMode = startAndDestinationModal.Mode;
    var defaultCredits = $("#defaultCredit").val();
    if (assetMode === "Simple" || assetMode === "DIAS_Simple") {
        startAndDestinationModal.StartLatitude = null;
        startAndDestinationModal.StartLongitude = null;
        startAndDestinationModal.DestinationLatitude = null;
        startAndDestinationModal.DestinationLongitude = null;
    }
    else {
        if (startLatitude === null || startLongitude === null) {
            msg = "Please provide Start Address </br>";
        }
        else {
            startAndDestinationModal.StartLatitude = startLatitude;
            startAndDestinationModal.StartLongitude = startLongitude;
        }
        if (destinationLatitude === null || destinationLongitude === null) {
            msg = msg + "Please provide Destination Address";
        }
        else {
            startAndDestinationModal.DestinationLatitude = destinationLatitude;
            startAndDestinationModal.DestinationLongitude = destinationLongitude;
        }

    }

    startAndDestinationModal.DefaultCredit = defaultCredits;
    if (msg !== "") {
        new PNotify({
            title: 'Info!',
            text: msg,
            type: 'info',
            delay: 2000
        });
        return;
    }
    else {
        $("#startDestinationModal").modal('hide');
        setQuestionStep();
        $("#questionDetailModal").modal('show');
        $("#questionModelBtn").show();
    }

}

function setQuestionStep() {
    $("#googleMap").removeClass("disabledMap");
    var assetMode = startAndDestinationModal.Mode;
    focusedInput = "questionAddress";

    $("#checkStartDest").hide();
    $("#checkStartDest").iCheck('uncheck');

    $("#DIASGrouppointAddressRadioBtnOnMap").hide();
    $("#DIASGrouppointAddressRadioBtn").iCheck('uncheck');

    $("#QuestionRadioBtnOnMap").show();
    $("#questionAddressRadioBtn").iCheck('check');

    if (assetMode === "Sequence") {
        $("#questionAddressDiv").show();
        $("#DIASGrouppointAddressDiv").hide();
        $("#btnAssociateQuesDiv").hide();

        $("#questionType").multiselect('select', "radio");
        $("#questionType").trigger('change');
        $("#questionType").multiselect('enable');

        
    }
    else if (assetMode === "Simple") {
        $("#questionAddressDiv").show();
        $("#DIASGrouppointAddressDiv").hide();
        $("#btnAssociateQuesDiv").hide();

        $("#questionType").multiselect('select', "radio");
        $("#questionType").trigger('change');
        $("#questionType").multiselect('enable');
    }
    else if (assetMode === "DIAS_Simple") {
        $("#questionAddressDiv").hide();
        $("#DIASGrouppointAddressDiv").show();
        $("#btnAssociateQuesDiv").hide();

        $("#questionType").multiselect('select', "likertScale");
        $("#questionType").trigger('change');
        $("#questionType").multiselect('disable');

        $("#DIASGrouppointAddressRadioBtnOnMap").show();
        $("#DIASGrouppointAddressRadioBtn").iCheck('check');

        $("#checkStartDest").hide();
        $("#checkStartDest").iCheck('uncheck');
        
        $("#QuestionRadioBtnOnMap").hide();
        $("#questionAddressRadioBtn").iCheck('check');

        focusedInput = "DIASGrouppointAddress";
    }
    else if (assetMode === "Decision") {
        $("#questionAddressDiv").show();
        $("#DIASGrouppointAddressDiv").hide();
        $("#btnAssociateQuesDiv").show();

        $("#questionType").multiselect('select', "radio");
        $("#questionType").multiselect('enable');
        $("#questionType").trigger('change');
        
    }
    LoadSensors();
}

// number validaiton for option's credits
$(document).on('keyup', '.optionsCredits', function (event) {

    if (((event.keyCode !== 13 || event.which !== 13)
        && (event.keyCode !== 8 || event.which !== 8)
        && (event.keyCode !== 9 || event.which !== 9))
        && (isNaN(parseInt($(this).val())))) {
        new PNotify({
            title: 'Info!',
            text: "Please enter a valid value",
            type: 'info',
            delay: 2000,
            nonblock: {
                nonblock: true
            }
        });
        $(this).val("");
    }
});
$(document).on("change", "#questionType", function () {
    resetOptionSection();
    optionCounter++;
    var textboxDivId = 'TextBoxDiv' + optionCounter;
    var newTextBoxDiv = $(document.createElement('div')).attr("id", textboxDivId);
    questionType = $("#questionType").val();

    if (questionType === "radio") {
        newTextBoxDiv.html(questionOptionHtml(optionCounter) + optionAddRemoveAndCreditHtml(optionCounter));
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");
    }
    else if (questionType === "checkbox") {
        
        newTextBoxDiv.html(questionOptionHtml(optionCounter) + optionAddRemoveButtonHtml(optionCounter));
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");

        $("#TextBoxesGroupCombination").html("");
        combinationCounter = 0; //resetting combination Counter to start from 1
        placeCombination("1", true, "");
    }

    else if (questionType === "likertScale") {
        newTextBoxDiv.html('<div class="slidecontainer" id="slidecontainer' + optionCounter + '">' +
            '<label style="width:15px;padding-left: 0px">1</label><input type="range" min="0" max="10" value="5" class="slider form-control-option-text" id="slider' + optionCounter + '"><label style="width:15px;margin-right:30px ;padding-left: 0px">10</label>' +
            '<input type="text" id="optQuestion' + optionCounter + '" class="optQuestion hide" >' +
            '<input type="number" id="optionsCredits' + optionCounter + '" min="1" placeholder="Credits" class="optionsCredits form-control form-control-option-credit">' +
            '</div>');
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");
    }

    else if (questionType === "textBox") {
        newTextBoxDiv.html('<textarea readonly id="textBox' + optionCounter + '" style="border: 1px solid lightgray;float: left;width: 70%;height: 34px;margin-top: 5px; margin-right: 10px;resize:none">User will enter response in a text field</textarea>' +
            '<input type="text" id="optQuestion' + optionCounter + '" class="optQuestion hide" >' +
            '<input type="number" id="optionsCredits' + optionCounter + '" min="1" placeholder="Credits" class="optionsCredits form-control form-control-option-credit">');
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");
    }
});

$(document).on("click", "#addButton", function () {
    assetMode = startAndDestinationModal.Mode;
    var questionType = $("#questionType").val();
    //********************************************
    //checking last entered option text
    //********************************************
    var v = $("#optQuestion" + optionCounter)
    var msg = "";
    if (v.val() === "" && ($("#questionType").val() !== "likertScale" && $("#questionType").val() !== "textBox"))
        msg = "Please enter option text."
    if (msg !== "") {
        new PNotify({
            title: 'Info!',
            text: msg,
            type: 'info',
            delay: 2000
        })
        return;
    }

    //checking number of options
    if (optionCounter >= 7) {
        if (questionType === "checkbox") {
            //placing combination for option#7
            var optsString = "";
            for (var i = 1; i <= optionCounter; i++) {
                optsString += i.toString();
            }
            combinationCounter = 0;
            $("#TextBoxesGroupCombination").html("")
            result = combinations(optsString);
            combinationCounter = 0;//resetting combination Counter to start from 1
            for (var k = 0; k < result.length; k++) {
                placeCombination(result[k], true, "");
            }
        }
        new PNotify({
            title: 'Info!',
            text: "Only 7 textboxes are allowed. ",
            type: 'info',
            delay: 2000
        })
        return false;
    }


    optionCounter++;
    var textboxDivId = 'TextBoxDiv' + optionCounter;
    var newTextBoxDiv = $(document.createElement('div')).attr("id", textboxDivId);
    var html = null;
    if (questionType === "radio") {
        newTextBoxDiv.html(questionOptionHtml(optionCounter) + optionAddRemoveAndCreditHtml(optionCounter));
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");
    }
    else if (questionType === "checkbox") {
        newTextBoxDiv.html(questionOptionHtml(optionCounter) + optionAddRemoveButtonHtml(optionCounter));
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");

        // placing combination
        var optsString = "";
        for (var i = 1; i <= optionCounter; i++) {
            optsString += i.toString();
        }
        $("#TextBoxesGroupCombination").html("")

        result = combinations(optsString);
        combinationCounter = 0; //resetting combination Counter to start from 1
        for (var k = 0; k < result.length; k++) {
            placeCombination(result[k], true, "");
        }

    }
});

$(document).on("click", ".classRemoveOption", function () {
    var id = $(this)[0].id.split("-");
    var optionId = id[id.length - 1];

    if (optionCounter <= 1) {
        new PNotify({
            title: 'Info!',
            text: "Its Mandatory to provide one option.",
            type: 'info',
            delay: 2000
        });
        return;
    }
    if (parseInt(optionId) === optionCounter) {
        $("#TextBoxDiv" + optionId).remove();
        optionCounter--;
    }
    else {
        $("#TextBoxDiv" + optionId).remove();
        for (var opt = (parseInt(optionId) + 1); opt <= optionCounter; opt++) {

            $("#TextBoxDiv" + opt).attr('id', "TextBoxDiv" + (opt - 1));
            $("#optQuestion" + opt).attr('id', "optQuestion" + (opt - 1));
            $("#removeOption-" + opt).attr('id', "removeOption-" + (opt - 1));
            $("#optionsCredits" + opt).attr('id', "optionsCredits" + (opt - 1));
        }
        optionCounter--;
    }


    if ($("#questionType").val() === "checkbox") {
        combinationCounter = 0; //resetting combination Counter to start from 1
        optionsString = ""
        for (var i = 1; i <= optionCounter; i++) {
            optionsString += i;
        }
        var result = combinations(optionsString);
        $("#TextBoxesGroupCombination").html("");
        for (var i = 0; i < result.length; i++)
            placeCombination(result[i], true, "");

    }
});

function addQuestionInList(qId) {
    var msg = "Please provide following information :</br>";
    if ($('#drpProject').val() == null || $('#drpProject').val() == "")
        msg += "=> Project</br>";

    if ($('#assetMode').val() == null || $('#assetMode').val() == "Sequence" || $('#assetMode').val() == "Decision") {
        if (startLatitude === null && startLongitude === null) {
            msg += "=> Start Address</br>"
        }
        if (destinationLatitude === null && destinationLongitude === null) {
            msg += "=> Destination Address</br>"
        }
    }

    if ($("#txtquestion").val() === "")
        msg += "=> Question</br>";

    if (questionLatitude === null && questionLongitude === null && !updateSts)
        msg += "=>Please mark question location on Map.</br>";

    if ($('#optQuestion' + optionCounter).val() === "" && ($("#questionType").val() !== "likertScale" && $("#questionType").val() !== "textBox"))
        msg += "=>Please enter option text.</br>";

    if ($("#txtVicinity").val() === "")
        msg += "=> Vicinity Value</br>";

    if ($("#txtTime").val() === "")
        msg += "=> Time Value</br>";

    if (msg != "Please provide following information :</br>") {
        new PNotify({
            title: 'Info!',
            text: msg,
            type: 'info',
            delay: 2000,
        });
        return;
    }

    for (var opt = 1; opt <= optionCounter; opt++) {
        var optionObj = {};
        optionObj.id = opt;
        optionObj.Credits = $("#optionsCredits" + opt).val();

        var credit = "Credits"
        if ($("#txtShowCredit").is(":checked") && $("#questionType").val() !== "checkbox") {
            if (!isNaN(parseInt(optionObj.Credits))) {
                if (parseInt(optionObj.Credits) == 1) {
                    credit = "Credit";
                }
                optionObj.Name = $("#optQuestion" + opt).val() + " (" + optionObj.Credits + " " + credit + ")";
            }
            else {
                if (parseInt($("#defaultCredit").val()) == 1) {
                    credit = "Credit";
                }
                optionObj.Name = $("#optQuestion" + opt).val() + " (" + $("#defaultCredit").val() + " " + credit + ")";
            }
        }
        else
            optionObj.Name = $("#optQuestion" + opt).val();
        lstOptions.push(optionObj);
    }

    if ($("#questionType").val() === "checkbox")
        AddCombination();
    var optionslistlength = lstOptions.length;

    var objQuestion = {};
    objQuestion.Question = $("#txtquestion").val();


    objQuestion.Latitude = questionLatitude;
    objQuestion.Longitude = questionLongitude;
    objQuestion.Type = $("#questionType").val();
    objQuestion.Time = $("#txtTime").val();
    objQuestion.Vicinity = $("#txtVicinity").val();
    objQuestion.Frequency = $("#txtFrequency").val();
    assetMode = $("#assetMode").val();
    if (assetMode === "Sequence") {
        objOption.NextQuestion = "Disablp;e";
        objQuestion.Sequence = $("#txtSequence").val();
    }
    if (assetMode === "Simple" || assetMode === "DIAS_Simple") {
        objOption.NextQuestion = "Disable";
        objQuestion.Sequence = "Disable";
    }
    if (objQuestion.Type === "radio" || objQuestion.Type === "likertScale" || objQuestion.Type === "textBox") {
        objQuestion.Combination = null;
    }
    if (objQuestion.Type === "checkbox") {

        objQuestion.Combination = lstCombinations;
    }
    objQuestion.Visibility = $("#txtVisibility").is(":checked");
    objQuestion.Mandatory = $("#txtMandatory").is(":checked");

    var Sensors = [];
    try {
        for (var i = 0; i < $("#dropdownSensors").val().map(String).length; i++) {
            var objSensor = {};
            objSensorId = objSensorId + 1;
            objSensor.id = objSensorId;
            objSensor.Name = $("#dropdownSensors").val()[i];
            Sensors.push(objSensor);
        }
    }
    catch (err) {
        new PNotify({
            title: 'Info!',
            text: "Please enter sensors.",
            type: 'info',
            delay: 2000
        });
        return;
    }
    objQuestion.Sensor = Sensors;
    objQuestion.Option = lstOptions;
    if (updateSts) {

        for (var l = 0; l < lstQuestionDataModel.length; l++) {
            var e = lstQuestionDataModel[l];
            if (e.id === parseInt(qId)) {
                objQuestionId = qId;
                objQuestion.id = objQuestionId;


                $("#main_panel_asset").animate({
                    backgroundColor: "#fff"
                }, 500);

                if ($("#assetMode").val() === "DIAS_Simple") {

                    questionAddressList[parseInt(qId) - 1] = $("#DIASGrouppointAddress").val();
                    lstQuestionDataModel[l] = objQuestion;
                    for (var n = 0; n < lstQuestionDataModel.length; n++) {
                        if ((questionOldLat === lstQuestionDataModel[n].Latitude) && (questionOldLng === lstQuestionDataModel[n].Longitude)) {
                            lstQuestionDataModel[n].Latitude = questionLatitude;
                            lstQuestionDataModel[n].Longitude = questionLongitude;
                            questionAddressList[(parseInt(lstQuestionDataModel[n].id) - 1)] = $("#DIASGrouppointAddress").val();
                        }
                    }

                }
                else {
                    lstQuestionDataModel[l].Latitude = questionLatitude;
                    lstQuestionDataModel[l].Longitude = questionLongitude;
                    questionAddressList[(parseInt(qId) - 1)] = $("#questionAddress").val();
                }
                
                var temp = new google.maps.Marker({
                    position: new google.maps.LatLng(questionLatitude, questionLongitude),
                    map: map,
                    icon: '/Content/img/blue-dot.png'
                });

                if ($("#assetMode").val() === "DIAS_Simple") {

                    for (var i = 0; i < questionLocationMarkers.length; i++) {

                        if (questionLocationMarkers[i].position.lat() == questionOldLat && questionLocationMarkers[i].position.lng() == questionOldLng) {
                            questionLocationMarkers[i].setPosition(
                                new google.maps.LatLng(
                                    questionLatitude,
                                    questionLongitude
                                )
                            );
                        }
                    }
                }
                else {
                    questionLocationMarkers[qId - 1] = temp

                }

                lstQuestionDataModel[l] = objQuestion;
                updateSts = false;
                var btnadd = '<button class="addquestion btn btn-default" id="btnaddquestion" onclick="addQuestionInList(0)">Add Question</button>';
                $("#addquestiondiv").html(btnadd);
            }
            
        }
    }
    else {
        objQuestionId = objQuestionId + 1;
        objQuestion.id = objQuestionId;
        lstQuestionDataModel.push(objQuestion);
        addedQuestionCount = addedQuestionCount + 1;
        
        var temp = new google.maps.Marker({
            position: new google.maps.LatLng(questionLatitude, questionLongitude),
            map: map,
            icon: '/Content/img/blue-dot.png'
        });
        if ($('#assetMode').val() === "DIAS_Simple") {
            questionAddressList.push($("#DIASGrouppointAddress").val());
        }
        else {
            questionAddressList.push($("#questionAddress").val());
        }
        questionLocationMarkers.push(temp);

    }

    new PNotify({
        title: 'Success',
        text: 'Question added Succesfully!',
        type: 'success',
        delay: 2000
    });
    $("#addedQuestiondiv_panel").animate({
        backgroundColor: "red"
    }, 500);


    $("#addedQuestiondiv_panel").animate({
        backgroundColor: "#fff"
    }, 2000);

    if (assetMode === "Decision" || assetMode === "Simple" || assetMode === "DIAS_Simple") {
        objQuestion.Sequence = "Disable";
    }

    // Disabling components
    $("#drpProject").multiselect('disable');
    $("#assetMode").attr("disabled", true);
    $("#startAddress").attr("disabled", true);
    $("#destinationAddress").attr("disabled", true);
    $("#btnStartAddress").attr("disabled", true);
    $("#btnDestinationAddress").attr("disabled", true);
    $("#defaultCredit").attr('disabled', true);
    //=====================

    LoadQuestions(lstQuestionDataModel);
    LoadInitialValues();

}

function LoadInitialValues() {
    $("#txtquestion").val('').blur();

    resetOptionSection();
    optionCounter++;
    if (startAndDestinationModal.Mode === "DIAS_Simple") {
        if (desiredLocationMarker !== null) {
            desiredLocationMarker.setMap(null);
            desiredLocationMarker = null;
        }
        var textboxDivId = 'TextBoxDiv' + optionCounter;
        var newTextBoxDiv = $(document.createElement('div')).attr("id", textboxDivId);
        newTextBoxDiv.html('<div class="slidecontainer" id="slidecontainer' + optionCounter + '">' +
            '<label style="width:15px;padding-left: 0px">1</label><input type="range" min="0" max="10" value="5" class="slider form-control-option-text" id="slider' + optionCounter + '"><label style="width:15px;margin-right:30px ;padding-left: 0px">10</label>' +
            '<input type="text" id="optQuestion' + optionCounter + '" class="optQuestion hide" >' +
            '<input type="number" id="optionsCredits' + optionCounter + '" min="1" placeholder="Credits" class="optionsCredits form-control form-control-option-credit">' +
            '</div>');
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");
        focusedInput = "DIASGrouppointAddress";
    }
    else {
        var val = $("#questionType").val();
        $("#questionType").multiselect('deselect', val);
        $("#questionType").multiselect("select", "radio");

        $("#questionAddress").val("");
        questionLatitude = null;
        questionLongitude = null;
        desiredLocationMarker.setMap(null);
        desiredLocationMarker = null;
        var textboxDivId = 'TextBoxDiv' + optionCounter;
        var newTextBoxDiv = $(document.createElement('div')).attr("id", textboxDivId);
        newTextBoxDiv.html(questionOptionHtml(optionCounter) + optionAddRemoveAndCreditHtml(optionCounter));
        newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
        newTextBoxDiv.appendTo("#TextBoxesGroup");
    }
}


function LoadSensors() {
    var options = [];
    options.push({ label: "Light", title: "Light", value: "Light" });
    options.push({ label: "Gyroscope", title: "Gyroscope", value: "Gyroscope" });
    options.push({ label: "Proximity", title: "Proximity", value: "Proximity" });
    options.push({ label: "Accelerometer", title: "Accelerometer", value: "Accelerometer" });
    options.push({ label: "Location", title: "Location", value: "Location", selected: "selected" });
    options.push({ label: "Noise", title: "Noise", value: "Noise" });
    $('#dropdownSensors').multiselect('dataprovider', options);

}

function associateQuestion() {

    $("#mainQuestionSelection").html("");
    boolAssociateQuestion = true;
    if (assetMode === "Decision") {
        if (lstQuestionDataModel.length >= 2) {
            if (boolAssociateQuestion === true) {
                loadOptionsInModal();
                $("#dgAssociateQuestion").css("display", "block");
                //.jQuery.noConflict();
                $("#dgAssociateQuestion").modal("show");
                $("#associationMsg").addClass("hide");
            }
        } else {
            new PNotify({
                title: 'Info!',
                text: "Please add at least two questions",
                type: 'info',
                delay: 2000
            });
        }
    }
}

$(document).on("change", "#mainQuestionSelection", function () {
    $("#associationMsg").addClass("hide");
    $("#btnSaveAssociateQuestion").removeClass("hide");

    $("#mbodyhoriAssociateOption").html("");
    mainQuestionSelectId = $(this).find("option:selected").attr('id');
    if (lstQuestionDataModel.length > 0) {
        var questionDropdwn = '';
        questionDropdwn = '<select id="associatedQuestionSelection" class="long"  style="border: 1px solid;">';
        questionDropdwn = questionDropdwn + '<option>Select Question</option>';
        for (var i = 0; i < lstQuestionDataModel.length; i++) {
            if (lstQuestionDataModel[i].id !== parseInt(mainQuestionSelectId)) {
                questionDropdwn = questionDropdwn + '<option id="' + lstQuestionDataModel[i].id + '">' + lstQuestionDataModel[i].Question + '</option>';
            }
        }

    } else {
        questionDropdwn = questionDropdwn + '<option> No question available.</option>';
    }
    questionDropdwn = questionDropdwn + '</select>';
    questionType = $("#questionType").val();
    for (var i = 0; i < lstQuestionDataModel.length; i++) {

        if (lstQuestionDataModel[i].id === parseInt(mainQuestionSelectId)) {
            var selectedQuestionOptions = lstQuestionDataModel[i].Option;
            lstOptionlen = selectedQuestionOptions.length

            if (lstQuestionDataModel[i].Combination === null) {
                if (lstQuestionDataModel[i].Type === "likertScale") {
                    var k = 0
                    for (; k < 10; k++) {
                        var populateOptions = '<div class="form-group" style="margin-top: 10px;">' +
                            '<div class="col-md-4">' +
                            '<label for="txtOption">value ' + (k + 1) + ':</label>' +
                            '</div>' +
                            '<div class="col-md-6" id=option' + (k + 1) + '>' +
                            questionDropdwn +
                            '</div>' +
                            '</div >';
                        $("#mbodyhoriAssociateOption").append(populateOptions);
                        if (selectedQuestionOptions[k] !== undefined && selectedQuestionOptions[k].NextQuestion !== undefined) {
                            $('#option' + (k + 1) + ' > select > option[id=' + selectedQuestionOptions[k].NextQuestion + ']').prop("selected", true);
                        }
                    }
                }

                else {
                    for (var k = 0; k < lstOptionlen; k++) {
                        var optionStr = "Option" + selectedQuestionOptions[k].id + " :";
                        if (lstQuestionDataModel[i].Type === "textBox") {
                            optionStr = "Option: TextBox";
                        }
                        var populateOptions = '<div class="form-group" style="margin-top: 10px;">' +
                            '<div class="col-md-4">' +
                            '<label for="txtOption">' + optionStr + selectedQuestionOptions[k].Name + '</label>' +
                            '</div>' +
                            '<div class="col-md-6" id=option' + selectedQuestionOptions[k].id + '>' +
                            questionDropdwn +
                            '</div>' +
                            '</div >';
                        $("#mbodyhoriAssociateOption").append(populateOptions);
                        if (selectedQuestionOptions[k].NextQuestion !== undefined) {

                            $('#option' + selectedQuestionOptions[k].id + ' > select > option[id=' + selectedQuestionOptions[k].NextQuestion + ']').prop("selected", true);
                        }
                    }
                }
            } else if (lstQuestionDataModel[i].Combination !== null) {
                var questionCombinations = lstQuestionDataModel[i].Combination;
                var combLength = questionCombinations.length;
                if (combLength !== 0) {
                    for (var k = 0; k < combLength; k++) {

                        var strcombination = '';
                        for (var j = 0; j < questionCombinations[k].Selected.length; j++) {
                            strcombination = strcombination + questionCombinations[k].Selected[j].Order;
                        }

                        var populateCombinations = '<div class="form-group" style="margin-top: 10px;">' +
                            '<div class="col-md-4">' +
                            '<label for="txtOption">Combination ' + questionCombinations[k].id + ': ' + strcombination + ' </label>' +
                            '</div>' +
                            '<div class="col-md-6" id=combination' + questionCombinations[k].id + '>' +
                            questionDropdwn +
                            '</div>' +
                            '</div >';
                        $("#mbodyhoriAssociateOption").append(populateCombinations);
                        if (questionCombinations[k].NextQuestion !== undefined) {
                            $('#combination' + questionCombinations[k].id + ' > select > option[id=' + questionCombinations[k].NextQuestion + ']').prop("selected", true);
                        }
                    }
                }
            }
        }
    }
});

$(document).on("click", "#btnSaveAssociateQuestion", function () {
    questionType = $("#questionType").val();
    boolAssociateQuestion = false;
    for (var i = 0; i < lstQuestionDataModel.length; i++) {
        if (lstQuestionDataModel[i].id === parseInt(mainQuestionSelectId)) {
            if (lstQuestionDataModel[i].Combination === null) {

                if (lstQuestionDataModel[i].Type === "likertScale") {
                    var name = lstQuestionDataModel[i].Option[0].Name;
                    var credit = lstQuestionDataModel[i].Option[0].Credits;
                    lstQuestionDataModel[i].Option[0].NextQuestion = $('#option' + 1 + ' option:selected').attr('id');
                    for (var k = 1; k < 10; k++) {
                        objOption = {};
                        objOption.id = k + 1;
                        objOption.Name = name;
                        objOption.Credits = credit;
                        objOption.NextQuestion = "";
                        lstQuestionDataModel[i].Option.push(objOption);
                        lstQuestionDataModel[i].Option[k].NextQuestion = $('#option' + (k + 1) + ' option:selected').attr('id');
                    }
                }
                else {
                    var selectedQuestionOptions = lstQuestionDataModel[i].Option;
                    lstOptionlen = selectedQuestionOptions.length
                    for (var k = 0; k < lstOptionlen; k++) {
                        lstQuestionDataModel[i].Option[k].NextQuestion = $('#option' + selectedQuestionOptions[k].id + ' option:selected').attr('id');
                    }
                }
            } else if (lstQuestionDataModel[i].Combination !== null) {
                var questionCombinations = lstQuestionDataModel[i].Combination;
                var combLength = questionCombinations.length;
                for (var k = 0; k < combLength; k++) {
                    lstQuestionDataModel[i].Combination[k].NextQuestion = $('#combination' + questionCombinations[k].id + ' option:selected').attr('id');
                }
            }
        }
    }
    $("#associationMsg").removeClass("hide");
    $("#btnSaveAssociateQuestion").addClass("hide");
    $("#mbodyhoriAssociateOption").html("");
});

function AddCombination() {
    
    lstCombinations = [];
    var SelectedOptions = [];
    var objSelectedOptionId = 0;
    objCombinationId = 0;

    optionsString = ""
    for (var i = 1; i <= optionCounter; i++) {
        optionsString += i;
    }
    var result = combinations(optionsString);
    optionsCombinations = [];

    for (var i = 0; i < result.length; i++) {
        if ($("#checkboxcombination-" + (i + 1)).is(':checked')) {
            var optionCombination = {};
            optionCombination.value = result[i];
            optionCombination.count = 1;
            optionsCombinations.push(optionCombination);

            var splittedSelectedCom = (optionCombination.value).split('');
            SelectedOptions = [];
            objSelectedOptionId = 0;
            for (var k = 0; k < splittedSelectedCom.length; k++) {
                var objSelectedOption = {};
                objSelectedOptionId = objSelectedOptionId + 1;
                objSelectedOption.id = objSelectedOptionId;
                objSelectedOption.Order = splittedSelectedCom[k];
                SelectedOptions.push(objSelectedOption);
            }
            objCombinationId = objCombinationId + 1;

            var txtAssociatedQuestionSelection = $('#associatedQuestionSelection' + (i + 1) + ' option:selected').attr('id');
            var txtCredits = $('#optionsCredits' + (i + 1)).val();
            var objCombination = {};
            objCombination.id = objCombinationId;
            objCombination.Selected = SelectedOptions;
            objCombination.NextQuestion = txtAssociatedQuestionSelection;
            objCombination.Credits = txtCredits;

            lstCombinations.push(objCombination);
        }
    }
    removeCombination = false;
}

function AddOption() {

    assetMode = startAndDestinationModal.Mode;
    objOptionId = objOptionId + 1;

    // If show credit check box checked then credit/s
    // will be appended with option so that the credit
    // show on android app.
    var txtOption = $('#optQuestion' + optionCounter).val();
    var txtCredits = $('#optionsCredits' + optionCounter).val();
    if ($("#txtShowCredit").is(":checked")) {
        var credit = "Credits"
        if (!isNaN(parseInt(txtCredits))) {
            if (parseInt(txtCredits) == 1) {
                credit = "Credit";
            }
            txtOption = txtOption + " (" + txtCredits + " " + credit + ")";
        }
        else {
            if (parseInt($("#defaultCredit").val()) == 1) {
                credit = "Credit";
            }
            txtOption = txtOption + " (" + $("#defaultCredit").val() + " " + credit + ")";
        }
    }

    // Setting Next Question Setting
    if (assetMode === "Sequence") {
        objOption.NextQuestion = "Disable";
    }
    if (assetMode === "Simple" || assetMode === "DIAS_Simple") {
        objOption.NextQuestion = "Disable";
    }
    if (assetMode === "Decision") {
        objOption.NextQuestion = $('#associatedQuestionSelection' + combinationCounter + ' option:selected').attr('id');
        if (objOption.NextQuestion === "") {
            objOption.NextQuestion = "Disable";
        }
    }

    if (questionType === "checkbox") {
        objOption.Credits = "Disable";
    }

    objOption = {};
    objOption.id = objOptionId;
    objOption.Name = txtOption;
    objOption.Credits = txtCredits;
    lstOptions.push(objOption);
    optionsString = "";
    for (var i = 0; i < lstOptions.length; i++) {
        optionsString = optionsString + lstOptions[i].id.toString();
    }

    var result = combinations(optionsString);

    optionsCombinations = [];
    for (var i = 0; i < result.length; i++) {
        var optionCombination = {};
        optionCombination.value = result[i];
        optionCombination.count = 1;
        optionsCombinations.push(optionCombination);
    }
    if (questionType === "checkbox") {
        $("#TextBoxesGroupCombination").html("");
        if (optionsCombinations.length !== 0) {
            $("#showCombinaions").html("");
            $("#showCombinaions").append('<label>Combinations: </label ><br />');
            var showcombinations = "";
            $("#TextBoxesGroupCombination").show();
            var lengthcom = optionsCombinations.length;
            for (var i = 0; i < lengthcom - 1; i++) {
                showcombinations = showcombinations + optionsCombinations[i].value + ' | ';
            }
            showcombinations = showcombinations + optionsCombinations[lengthcom - 1].value;
            combinationCounter = 0; //resetting combination Counter to start from 1

            $("#showCombinaions").show();
            $("#showCombinaions").append(showcombinations);
        }
    } else {
        $("#showCombinaions").hide();
    }
    optionsAdded += 1;
}

function placeCombination(k, isSelected, creditVal) {
    combinationCounter++;
    var textboxDivCombinationId = 'TextBoxDivCombinaion' + combinationCounter;
    var newTextBoxDivCombinaion = $(document.createElement('div')).attr("id", textboxDivCombinationId);
    var combinationtags = '<label id="lblCombination-' + combinationCounter + '" class="lblCombinaitonTxt">' + k + '</label>' +
        '<input type="checkbox" onchange="checkCombinations(this);" id="checkboxcombination-' + combinationCounter + '" style="margin-right: 0;">' +
        '<input type="number" min="1" id="optionsCredits' + combinationCounter + '" placeholder="Credits" value="' + creditVal + '" class="optionsCredits form-control form-control-combination-credit">'
    newTextBoxDivCombinaion.html(combinationtags);
    newTextBoxDivCombinaion.addClass("textboxDiv col-lg-4 col-md-4 col-sm-4 col-xs-4");
    $("#TextBoxesGroupCombination").append(newTextBoxDivCombinaion);
    if (isSelected)
        $('#checkboxcombination-' + combinationCounter).prop('checked', true);

}

function checkCombinations(event) {

    var id = event.id;
    if (id.includes("checkboxcombination")) {
        var isChecked = event.checked;
        var splitedcombinId = id.split('-')[1];
        var lblvalue = $("#lblCombination-" + splitedcombinId).text();
        for (var i = 0; i < optionsCombinations.length; i++) {
            if (lblvalue === optionsCombinations[i].value) {
                if (isChecked === true) {
                    optionsCombinations[i].count = 1;
                } else if (isChecked === false) {
                    optionsCombinations[i].count = 0;
                    $('#associatedQuestionSelection' + combinationCounter).html('');
                    $('#optionsCredits' + combinationCounter).html('');
                }
            }
        }
    }
}

function LoadQuestions(lstQuestions) {
    $.fn.dataTable.ext.errMode = "throw";
    $("#tbQuestions").empty();
    $("#tbQuestions").on("error.dt", function (e, settings, techNote, message) {
        console.log(message);
    }).DataTable({
        language: {
            emptyTable: "No Questions available."
        },
        scrollY: "250px",
        data: lstQuestions,
        autoWidth: false,
        info: false,
        lengthChange: false,
        destroy: true,
        order: [[0, "asc"]],
        columnDefs: [
            {
                className: "coloumnWrap",
                targets: 0,
                width: '5px',
            },
            {
                className: "coloumnWrap",
                targets: 1,
                width: '20px'
            },
            {
                className: "coloumnWrap",
                targets: 2,
                width: '20px'
            },
            {
                className: "coloumnWrap",
                targets: 3,
                width: '10px'
            },
            {
                className: "coloumnWrap",
                targets: 4,
                width: '10px'
            },

            {
                className: "coloumnWrap",
                targets: 5,
                width: '20px'
            },
            {
                className: "icons-header align-vertical-top",
                targets: 6,
                width: '15px'
            }
        ],
        columns: [
            {
                data: "id",
                searchable: false,
                title: "Id",
                width: '5%'
            },
            {
                data: "Question",
                searchable: false,
                title: "Question",
                width: '20px'
            },
            {
                data: "Sensor",
                searchable: false,
                render: function (data, type, full) {
                    var s = "";
                    for (var i = 0; i < full.Sensor.length; i++) {
                        s = s + full.Sensor[i].Name + "</br>";
                    }
                    return s
                },
                title: "Sensors",
                width: '20px'
            },
            {
                data: "Type",
                searchable: false,
                title: "Type",
                width: '10px'
            },
            {

                width: "10px",
                title: "Location",
                render: function (data, type, full) {
                    if (questionAddressList.length > 0)
                        return questionAddressList[full.id - 1];
                    else
                        return "";
                }
            },
            {
                data: "Option",
                searchable: false,
                render: function (data, type, full) {
                    var o = "";

                    for (var i = 0; i < full.Option.length; i++) {
                        o = o + full.Option[i].Name + "</br>";
                    }
                    return o
                },
                title: "Options",
                width: '20px'
            },
            {
                data: "id",
                sortable: false,
                render: function (data, type, full) {
                    var icons = '';
                    icons = icons +
                        '<img class="btn btnQuestionEdit icons" id="editQuestionId" src="../Content/img/ico-edit.png" title="Edit">' +
                        '<span class="grid-icon-separator">' +
                        '<img class="btn btn-xxs deletebtn btnQuestionDelete icons" type="button" id="deleteQuestionId" src="../Content/img/ico-delete.png" title="Delete">'
                    return icons;
                },
                title: "&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp",
                width: '15px'
            }
        ],
        initComplete:
            function () {
                $(this.api().table().header()).css({ 'background-color': "#ECEBEB", 'color': "#000", 'width': '100%' });
            }
    });
    $('.dataTables_scrollHeadInner,.dataTables_scrollHeadInner table').width("100%");

}

function resetOptionSection() {
    
    optionCounter = 0;
    combinationCounter = 0;
    assetMode = $("#assetMode").val();

    objCombinationId = 0;
    objCombination = {};
    lstCombinations = [];

    objOptionId = 0;
    objOption = {};
    lstOptions = [];

    $("#TextBoxesGroup").html("");
    $("#TextBoxesGroupCombination").html("");
}

$(document).on('click', "#deleteQuestionId", function () {
    if (updateSts) {
        new PNotify({
            title: 'Info!',
            text: "Please update selected question first",
            type: 'info',
            delay: 2000
        });
    }
    else {
        var questionId = $(this).closest('tr').find("td:nth-child(1)").text();

        lstQuestionDataModel = $.grep(lstQuestionDataModel, function (e) {

            if (e.id === parseInt(questionId)) {
                return false;
            } else {
                return true;
            }
        });

        if (lstQuestionDataModel.length == 0)
            refreshPage();
        LoadQuestions(lstQuestionDataModel);
        questionLocationMarkers[questionId - 1].setMap(null);
        questionLocationMarkers[questionId - 1] = null;
    }
});

$(document).on('click', "#editQuestionId", function () {
    if (updateSts) {
        new PNotify({
            title: 'Info!',
            text: "Please update selected question first",
            type: 'info',
            delay: 2000
        });
    }
    else {
        var questionId = $(this).closest('tr').find("td:nth-child(1)").text();
        $(this).addClass("hide");
        editedQuestionId = questionId;
        updateSts = true;
        $("#main_panel_asset").animate({
            backgroundColor: "#EDEDED"
        }, 500);
        LoadQuestionById(questionId);
    }
});

function LoadQuestionById(questionId) {

    assetMode = startAndDestinationModal.Mode;
    for (var i = 0; i < questionLocationMarkers.length; i++) {
        if (questionId == (i + 1)) {
            questionOldLat = questionLocationMarkers[i].position.lat();
            questionOldLng = questionLocationMarkers[i].position.lng();
            questionLatitude = questionLocationMarkers[i].position.lat();
            questionLongitude = questionLocationMarkers[i].position.lng();
            desiredLocationMarker = questionLocationMarkers[i];
        }
    }

    for (var l = 0; l < lstQuestionDataModel.length; l++) {
        var e = lstQuestionDataModel[l];
        if (e.id === parseInt(questionId)) {

            $("#txtquestion").val(e.Question);
            if (assetMode === "DIAS_Simple") {
                $("#DIASGrouppointAddress").val(questionAddressList[(parseInt(questionId)) - 1]);
            }
            else {
                $("#questionAddress").val(questionAddressList[(parseInt(questionId)) - 1]);
            }
            $("#questionLatitude").val(e.Latitude);
            $("#questionLongitude").val(e.Longitude);
            var sensors = [];
            for (var n = 0; n < e.Sensor.length; n++) {
                sensors.push(e.Sensor[n].Name);
            }

            $("#questionAddress").val();
            $("#dropdownSensors").multiselect('select', sensors);
            $("#questionType").multiselect('select', e.Type);
            $("#txtTime").val(e.Time);
            $("#txtVicinity").val(e.Vicinity);
            $("#txtFrequency").multiselect('select', e.Frequency);
            if (e.Sequence !== "Disbale") {
                $("#lblSequence").show();
                $("#txtSequence").show();
                $("#txtSequence").val(e.Sequence);
            }

            if (e.Visibility === true) {
                if ($("#txtVisibility").prop("checked") == false)
                    $("#txtVisibility").trigger('click').attr("checked", "checked");
            }
            else {
                if ($("#txtVisibility").prop("checked") == true)
                    $("#txtVisibility").trigger('click').removeAttr("checked");
            }


            if (e.Mandatory === true) {
                if ($("#txtMandatory").prop("checked") == false)
                    $("#txtMandatory").trigger('click').attr("checked", "checked");
            }
            else {
                if ($("#txtMandatory").prop("checked") == true)
                    $("#txtMandatory").trigger('click').removeAttr("checked");
            }


            //*******************************
            //loading options
            //*******************************
            if (e.Type === "radio") {
                var maindiv = $("#TextBoxesGroup");
                maindiv.html("");
                optionCounter = e.Option.length;
                for (var i = 0; i < e.Option.length; i++) {
                    var textboxDivId = 'TextBoxDiv' + e.Option[i].id;
                    var div = $(document.createElement('div')).attr('id', textboxDivId);
                    div.html(questionOptionHtml(e.Option[i].id) + optionAddRemoveAndCreditHtml(e.Option[i].id));
                    maindiv.append(div);

                    var credit = "Credits";

                    if (parseInt(e.Option[i].Credits) === 1)
                        credit = "Credit";
                    if (isNaN(parseInt(e.Option[i].Credits))) {
                        $("#optQuestion" + e.Option[i].id).val(e.Option[i].Name.replace("(" + $("#defaultCredit").val() + " " + credit + ")", ""));
                    }
                    else {
                        $("#optQuestion" + e.Option[i].id).val(e.Option[i].Name.replace("(" + e.Option[i].Credits + " " + credit + ")", ""));
                    }
                    $("#optionsCredits" + e.Option[i].id).val(e.Option[i].Credits);
                    div.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
                }
            }


            if (e.Type === "checkbox") {
                var maindiv = $("#TextBoxesGroup");
                maindiv.html("");
                optionCounter = e.Option.length;
                for (var i = 0; i < e.Option.length; i++) {
                    var textboxDivId = 'TextBoxDiv' + e.Option[i].id;
                    var div = $(document.createElement('div')).attr('id', textboxDivId);
                    div.html(questionOptionHtml(e.Option[i].id) + optionAddRemoveButtonHtml(e.Option[i].id));
                    maindiv.append(div);
                    $("#optQuestion" + e.Option[i].id).val(e.Option[i].Name)
                    div.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
                }
                var optsString = "";
                for (var i = 1; i <= optionCounter; i++) {
                    optsString += i.toString();
                }
                $("#TextBoxesGroupCombination").html("")

                lstlabelsCreditsOfSelectedComb = gettingSelectedcombinationLabel(e.Combination);
                result = combinations(optsString);
                combinationCounter = 0; //resetting combination Counter to start from 1
                for (var k = 0; k < result.length; k++) {
                    
                    var index = $.inArray(result[k], lstlabelsCreditsOfSelectedComb.labels);
                    if (index !== -1) {
                        placeCombination(result[k], true, lstlabelsCreditsOfSelectedComb.credits[index]);
                    }
                    else
                        placeCombination(result[k], false, "");
                }
            }


            else if (e.Type === "likertScale") {
                var textboxDivId = 'TextBoxDiv' + optionCounter;
                var newTextBoxDiv = $(document.createElement('div')).attr("id", textboxDivId);

                var val = "";
                if (isNaN(e.Option[0].Credits))
                    val = $("#defaultCredit").val();
                else
                    val = e.Option[0].Credits
                newTextBoxDiv.html('<div class="slidecontainer" id="slidecontainer' + optionCounter + '">' +
                    '<label style="width:15px;padding-left: 0px">1</label><input type="range" min="0" max="10" value="5" class="slider form-control-option-text" id="slider' + optionCounter + '"><label style="width:15px;margin-right:30px ;padding-left: 0px">10</label>' +
                    '<input type="text" id="optQuestion' + optionCounter + '" class="optQuestion hide" >' +
                    '<input type="number" min="1" id="optionsCredits' + optionCounter + '" placeholder="Credits" class="optionsCredits form-control form-control-option-credit" value="' + val + '">' +
                    '</div>');
                newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
                $("#TextBoxesGroup").html("");
                newTextBoxDiv.appendTo("#TextBoxesGroup");
            }

            else if (e.Type === "textBox") {
                var textboxDivId = 'TextBoxDiv' + optionCounter;
                var newTextBoxDiv = $(document.createElement('div')).attr("id", textboxDivId);
                var val = "";
                if (isNaN(e.Option[0].Credits))
                    val = $("#defaultCredit").val();
                else
                    val = e.Option[0].Credits
                newTextBoxDiv.html('<textarea readonly id="textBox' + optionCounter + '" style="border: 1px solid lightgray;float: left;width: 70%;height: 34px;margin-top: 5px; margin-right: 10px;resize:none">User will enter response in a text field</textarea>' +
                    '<input type="text" id="optQuestion' + optionCounter + '" class="optQuestion hide" >' +
                    '<input type="number" id="optionsCredits' + optionCounter + '" min="1" placeholder="Credits" class="optionsCredits form-control form-control-option-credit" value="' + val + '">');
                newTextBoxDiv.addClass("textboxDiv col-lg-6 col-md-6 col-sm-6 col-xs-6");
                $("#TextBoxesGroup").html("");
                newTextBoxDiv.appendTo("#TextBoxesGroup");
            }

            $("#btnaddquestion").hide();
            var btnupdate = '<button class="addquestion btn btn-default" id="btnaddquestion" onclick="addQuestionInList(' + questionId + ')">Update Question</button>';
            $("#addquestiondiv").html(btnupdate);
        }
    }
}
function gettingSelectedcombinationLabel(lstSelectedCombination) {
    
    var lstSelectedCombinationLabel = [];
    var lstSelectedCombinationCredit = [];
    var lstSelectedCombinationLabelAndCredit = {};
    for (var i = 0; i < lstSelectedCombination.length; i++) {
        var label = "";
        for (var j = 0; j < lstSelectedCombination[i].Selected.length; j++) {
            label += lstSelectedCombination[i].Selected[j].Order;
        }
        lstSelectedCombinationCredit.push(lstSelectedCombination[i].Credits)
        lstSelectedCombinationLabel.push(label);
    }
    lstSelectedCombinationLabelAndCredit.labels = lstSelectedCombinationLabel;
    lstSelectedCombinationLabelAndCredit.credits = lstSelectedCombinationCredit;
    return lstSelectedCombinationLabelAndCredit;
}
function submitFormevent(event) {
    console.log("Function called");
    if (addedQuestionCount === 0) {
        new PNotify({
            title: 'Info!',
            text: "Please Click on Add Question button to proceed",
            type: 'info',
            delay: 2000
        });
        return;
    }
    $("#btnsubmit").addClass('buttonDisabled');
    $("#btnDownloadXMLFile").removeClass('buttonDisabled');
    if (lstQuestionDataModel === null) {
        new PNotify({
            title: 'Info!',
            text: "Please add atleast one question.",
            type: 'info',
            delay: 2000
        });
        return;
    }
    objMainModel.StartAndDestinationModel = startAndDestinationModal;
    objMainModel.SampleDataModel = lstQuestionDataModel;
    GenerateXMLCall();
}

function submitForHive(event) {
    if (addedQuestionCount === 0) {
        new PNotify({
            title: 'Info!',
            text: "Please Click on Add Question button to proceed",
            type: 'info',
            delay: 2000
        });
        return;
    }

    if (lstQuestionDataModel === null) {
        new PNotify({
            title: 'Info!',
            text: "Please add atleast one question.",
            type: 'info',
            delay: 2000
        });
        return;
    }
    objMainModel.StartAndDestinationModel = startAndDestinationModal;
    objMainModel.SampleDataModel = lstQuestionDataModel;

    GenerateHiveCall();
}

function GenerateHiveCall() {
    var projectId = document.getElementById('drpProject').value;
    if (projectId == null || projectId == "") {
        new PNotify({
            title: 'Info!',
            text: "Please select Project before submitting to Hive.",
            type: 'info',
            delay: 2000
        })
        return;
    }
    var mainModel = {};
    mainModel.ProjectId = projectId;
    mainModel.QuestionsModel = objMainModel;

    $.ajax({
        type: "POST",
        url: generateHiveCall,
        data: JSON.stringify(mainModel),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        async: true,
        beforeSend: function () {
            showLoading();
        },
        complete: function () {
            hideLoading()
        },
        success: function (data, textStatus, jqXHR) {
            if (data === "LogOut") {
                new PNotify({
                    title: 'Info!',
                    text: "Please Login to Submit! ",
                    type: 'info',
                    delay: 2000,

                    after_close: function () {
                        if (result.value) {
                            var url = "/Home/Index";
                            window.location.href = url;
                        }
                    },
                });
            }
            else {
                ajaxreturncode = jqXHR.status;

                var username = '<%= Session["UserName"] %>';
                new PNotify({
                    title: "Success!",
                    text: 'Asset Created Successfully!',
                    type: 'success',
                    delay: 1500,
                    after_close: function () {
                        $("#formAnchor")[0].click();
                    },

                });

            }
        },
        error: function (data, textStatus, jqXHR) {
            ajaxreturncode = jqXHR.status;
            new PNotify({
                title: 'Error!',
                text: textStatus,
                type: 'error',
                delay: 2000
            });
        }
    });

}

function GenerateXMLCall() {
    var mainModel = {};
    var pId = document.getElementById('drpProject').value;

    mainModel.ProjectId = pId;
    mainModel.QuestionsModel = objMainModel;

    $.ajax({
        type: "POST",
        url: generateXMLFile,
        data: JSON.stringify(mainModel),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        async: true,
        complete: function () {
            new PNotify({
                title: 'Success!',
                text: "Json file is created.",
                type: 'success',
                delay: 2000

            });
        },
        success: function (data, textStatus, jqXHR) {
            ajaxreturncode = jqXHR.status;
        },
        error: function (data, textStatus, jqXHR) {
            ajaxreturncode = jqXHR.status;
        }
    });

}

function GetDownloadedFileCall() {
    return $.ajax({
        url: getdownloadedFile,
        async: false
    });
}

function DeleteDownloadedFileCall() {
    return $.ajax({
        type: "DELETE",
        url: deletedownloadedFile,
        async: true
    });
}

function downloadFile() {
    GetDownloadedFileCall().success(function (response) {
        download("JSONFile.json", response);
    });
    DeleteDownloadedFileCall().success(function (response) {
    });
    $("#btnsubmit").removeClass("buttonDisabled");
    $("#btnDownloadXMLFile").addClass("buttonDisabled");
}

function download(filename, text) {
    var element = document.createElement('a');
    element.setAttribute('href', 'data:application/json,' + encodeURIComponent(text));
    element.setAttribute('download', filename);
    element.style.display = 'none';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
}

function customException(m) {
    throw m;
}

function uploadFile() {

    var file = $("#assetFile")[0].files[0];

    if (file) {
        var fileType = file.name.substr(file.name.lastIndexOf("."))
        if (fileType !== ".json") {
            new PNotify({
                title: 'Info!',
                text: "Please select json file.",
                type: 'info',
                delay: 2000

            });
            return;
        }

        var mainModelFromFile = {};

        var v = "";
        // create reader
        var reader = new FileReader();
        reader.readAsText(file);
        reader.onload = function (e) {
            v = JSON.parse(e.target.result);

        }
        setTimeout(function () {
            try {
                var objMainModelFromFile = {};
                //Parsing Project Id and verifying the given 
                //project id is belongs to the currently Logged in User.
                var pId = "";
                var msg = "";
                v.ProjectId !== undefined ? (pId = v.ProjectId) : msg = msg + '=>"ProjectId" key not found</br>';
                mainModelFromFile.ProjectId = "";
                if (pId !== "") {
                    for (var i = 0; i < projectIds.length; i++) {
                        if (projectIds[i].id === pId) {
                            mainModelFromFile.ProjectId = pId;
                            break;
                        }
                    }
                }


                //Parsing QuestionsModel
                var questionsModel = {};
                v.QuestionsModel !== undefined ? questionsModel = v.QuestionsModel : msg = msg + '=>"QuestionsModel" key not found.</br>'

                // Parsing Start and Destination Model.
                var startAndDestinationModel = null;
                questionsModel.StartAndDestinationModel !== undefined ? startAndDestinationModel = questionsModel.StartAndDestinationModel : msg = msg + '=>"StartAndDestinaitonModel" key inside QuestionsModel not found.</br>';
                if (startAndDestinationModel !== null) {
                    var mode = null;
                    (startAndDestinationModel.Mode !== undefined) ? mode = startAndDestinationModel.Mode : msg = msg + '=>"Mode" key inside StartAndDestinaitonModel not found.</br>';
                    if (mode !== null && (mode === "Simple" || mode === "DIAS_Simple")) {
                        (startAndDestinationModel.StartLatitude !== undefined && startAndDestinationModel.StartLatitude !== null) ? msg = msg + '=>"StartLatitude" key inside StartAndDestinaitonModel must be null.</br>' : startAndDestinationModel.StartLatitude = null;
                        (startAndDestinationModel.StartLongitude !== undefined && startAndDestinationModel.StartLongitude !== null) ? msg = msg + '=>"StartLongitude" key inside StartAndDestinaitonModel must be null.</br>' : startAndDestinationModel.StartLongitude = null;
                        (startAndDestinationModel.DestinationLatitude !== undefined && startAndDestinationModel.DestinationLatitude !== null) ? msg = msg + '=>"DestinationLatitude" key inside StartAndDestinaitonModel must be null.</br>' : startAndDestinationModel.DestinationLatitude = null;
                        (startAndDestinationModel.DestinationLongitude !== undefined && startAndDestinationModel.DestinationLongitude !== null) ? msg = msg + '=>"DestinationLongitude" key inside StartAndDestinaitonModel must be null</br>' : startAndDestinationModel.DestinationLongitude = null;
                    }
                    else {
                        (startAndDestinationModel.StartLatitude !== undefined) ? startAndDestinationModel.StartLatitude : msg = msg + '=>"StartLatitude" key inside StartAndDestinaitonModel not found.</br>';
                        (startAndDestinationModel.StartLongitude !== undefined) ? startAndDestinationModel.StartLongitude : msg = msg + '=>"StartLongitude" key inside StartAndDestinaitonModel not found.</br>';
                        (startAndDestinationModel.DestinationLatitude !== undefined) ? startAndDestinationModel.DestinationLatitude : msg = msg + '=>"DestinationLatitude" key inside StartAndDestinaitonModel not found.</br>';
                        (startAndDestinationModel.DestinationLongitude !== undefined) ? startAndDestinationModel.DestinationLongitude : msg = msg + '=>"DestinationLongitude" key inside StartAndDestinaitonModel not found.</br>';
                    }
                    (startAndDestinationModel.DefaultCredit !== undefined) ? startAndDestinationModel.defaultCredits = startAndDestinationModel.DefaultCredit : msg = msg + '=>"DefaultCredit" key inside StartAndDestinaitonModel not found.</br>';

                }
                objMainModelFromFile.StartAndDestinationModel = startAndDestinationModel;


                //Parsing SamleDataModel
                var sampleDataModel = null;
                questionsModel.SampleDataModel !== undefined ? (sampleDataModel = questionsModel.SampleDataModel) : msg = msg + '=>"SampleDataModel" key inside QuestionsModel not found.</br>';
                if (sampleDataModel !== null) {
                    var lstQuestionDataModelFromFile = [];
                    for (var i = 0; i < sampleDataModel.length; i++) {

                        var objQuestionFromFile = {};

                        objQuestionFromFile.id = (sampleDataModel[i].id !== undefined) ? sampleDataModel[i].id : msg = msg + '=>"id" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        objQuestionFromFile.Question = (sampleDataModel[i].Question !== undefined) ? sampleDataModel[i].Question : msg = msg + '=>"Question" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        objQuestionFromFile.Type = (sampleDataModel[i].Type !== undefined) ? sampleDataModel[i].Type : msg = msg + '=>"Type" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        objQuestionFromFile.Latitude = (sampleDataModel[i].Latitude !== undefined) ? sampleDataModel[i].Latitude : msg = msg + '=>"Latitude" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        objQuestionFromFile.Longitude = (sampleDataModel[i].Longitude !== undefined) ? sampleDataModel[i].Longitude : msg = msg + '=>"Longitude" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        objQuestionFromFile.Sensor = null;
                        var sensorValues = null;
                        (sampleDataModel[i].Sensor !== undefined) ? sensorValues = sampleDataModel[i].Sensor : msg = msg + '=>"Sensor" key inside SampleDataModel question object no. ' + i + ' not found.</br>';

                        if (sensorValues !== null) {
                            var sensorLst = [];
                            for (var s = 0; s < sensorValues.length; s++) {
                                var sensorObj = {};
                                sensorObj.id = null;
                                (sensorValues[s].id !== undefined) ? sensorObj.id = sensorValues[s].id : msg = msg + '=>"id" key inside SampleDataModel question object no. ' + i + ' sensor object no. ' + s + ' not found.</br>';
                                sensorObj.Name = null;
                                (sensorValues[s].Name !== undefined) ? sensorObj.Name = sensorValues[s].Name : msg = msg + '=>"Name" key inside SampleDataModel question object no. ' + i + ' sensor object no. ' + s + ' not found.</br>';
                                sensorLst.push(sensorObj);
                            }
                            objQuestionFromFile.Sensor = sensorLst;
                        }

                        objQuestionFromFile.Time = (sampleDataModel[i].Time !== undefined) ? sampleDataModel[i].Time : msg + '=>"Time" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        objQuestionFromFile.Frequency = (sampleDataModel[i].Frequency !== undefined) ? sampleDataModel[i].Frequency : msg + '=>"Frequency" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        if (objMainModelFromFile.StartAndDestinationModel.Mode === "Sequence")
                            objQuestionFromFile.Sequence = (sampleDataModel[i].Sequence !== undefined) ? sampleDataModel[i].Sequence : msg + '=>"Sequence" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        else
                            objQuestionFromFile.Sequence = "Disable";
                        objQuestionFromFile.Visibility = (sampleDataModel[i].Visibility !== undefined) ? sampleDataModel[i].Visibility : msg + '=>"Visibility" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        objQuestionFromFile.Mandatory = (sampleDataModel[i].Mandatory !== undefined) ? sampleDataModel[i].Mandatory : msg + '=>"Mandatory" key inside SampleDataModel question object no. ' + i + ' not found.</br>';

                        //Parsing Options Data
                        var optionValues = null;
                        objQuestionFromFile.Option = null;

                        (sampleDataModel[i].Option !== undefined) ? optionValues = sampleDataModel[i].Option : msg = msg + '=>"Option" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        if (optionValues !== null) {
                            var optionLst = [];
                            for (var o = 0; o < optionValues.length; o++) {
                                var optionObj = {};
                                optionObj.id = null;
                                (optionValues[o].id !== undefined) ? optionObj.id = optionValues[o].id : msg = msg + '=>"id" key inside SampleDataModel question object no. ' + i + ' option object no. ' + o + ' not found.</br>';
                                optionObj.Name = null;
                                (optionValues[o].Name !== undefined) ? optionObj.Name = optionValues[o].Name : msg = msg + '=>"Name" key inside SampleDataModel question object no. ' + i + ' option object no. ' + o + ' not found.</br>';
                                optionObj.NextQuestion = null;
                                if (objMainModelFromFile.StartAndDestinationModel.Mode === "Sequence" || objMainModelFromFile.StartAndDestinationModel.Mode === "Simple" || objMainModelFromFile.StartAndDestinationModel.Mode === "DIAS_Simple")
                                    optionObj.NextQuestion = "Disable";
                                else if (objMainModelFromFile.StartAndDestinationModel.Mode === "Decision" && objQuestionFromFile.Type === "checkbox")
                                    optionObj.NextQuestion = null;
                                else
                                    (optionValues[o].NextQuestion !== undefined) ? optionObj.NextQuestion = optionValues[o].NextQuestion : msg = msg + '=>"NextQuestion" key inside SampleDataModel question object no. ' + i + ' option object no. ' + o + ' not found.</br>';

                                optionObj.Credits = null;
                                if (objQuestionFromFile.Type === "checkbox")
                                    optionObj.Credits = "Disable";
                                else
                                    (optionValues[o].Credits !== undefined) ? optionObj.Credits = optionValues[o].Credits : msg = msg + '=>"Credits" key inside SampleDataModel question object no. ' + i + ' option object no. ' + o + ' not found.</br>';
                                optionLst.push(optionObj);
                            }
                            objQuestionFromFile.Option = optionLst;
                        }

                        //Parsing Combination Data
                        objQuestionFromFile.Combination = null;
                        if (objQuestionFromFile.Type !== undefined && objQuestionFromFile.Type === "checkbox") {
                            var combinationValues = null;
                            (sampleDataModel[i].Combination !== undefined) ? combinationValues = sampleDataModel[i].Combination : msg = msg + '=>"Combination" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                            if (combinationValues !== null) {
                                var combinationLst = [];
                                for (var c = 0; c < combinationValues.length; c++) {
                                    var combinationObj = {};
                                    combinationObj.id = null;
                                    (combinationValues[c].id !== undefined) ? combinationObj.id = combinationValues[c].id : msg = msg + '=>"id" key inside SampleDataModel question object no. ' + i + ' combination object no. ' + c + ' not found.</br>';

                                    combinationObj.NextQuestion = null;
                                    if (objMainModelFromFile.StartAndDestinationModel.Mode === "Decision")
                                        (combinationValues[c].NextQuestion !== undefined) ? combinationObj.NextQuestion = combinationValues[c].NextQuestion : msg = msg + '=>"NextQuestion" key inside SampleDataModel question object no. ' + i + ' combination object no. ' + c + ' not found.</br>';
                                    combinationObj.Credits = null;
                                    (combinationValues[c].Credits !== undefined) ? combinationObj.Credits = combinationValues[c].Credits : msg = msg + '=>"Credits" key inside SampleDataModel question object no. ' + i + ' combination object no. ' + c + ' not found.</br>';

                                    combinationObj.Selected = null;
                                    var selectedValues = null;
                                    (combinationValues[c].Selected !== undefined) ? selectedValues = combinationValues[c].Selected : msg = msg + '=>"Selected" key inside SampleDataModel question object no. ' + i + ' combination object no. ' + c + ' not found.</br>';
                                    if (selectedValues !== null) {
                                        var selectedValuesLst = [];
                                        for (var sv = 0; sv < selectedValues.length; sv++) {
                                            var selectedValuesObj = {};
                                            selectedValuesObj.id = null;
                                            (selectedValues[sv].id !== undefined) ? selectedValuesObj.id = selectedValues[sv].id : msg = msg + '=>"id" key inside SampleDataModel question object no. ' + i + ' combination object no. ' + c + ' selected object no. ' + sv + ' not found.</br>';
                                            selectedValuesObj.Order = null;
                                            (selectedValues[sv].Order !== undefined) ? selectedValuesObj.Order = selectedValues[sv].Order : msg = msg + '=>"Order" key inside SampleDataModel question object no. ' + i + ' combination object no. ' + c + ' selected object no. ' + sv + ' not found.</br>';
                                            selectedValuesLst.push(selectedValuesObj);
                                        }
                                        combinationObj.Selected = selectedValuesLst;
                                    }
                                    combinationLst.push(combinationObj);
                                }
                                objQuestionFromFile.Combination = combinationLst;

                            }
                        }

                        objQuestionFromFile.Vicinity = (sampleDataModel[i].Vicinity !== undefined) ? sampleDataModel[i].Vicinity : msg + '=>"Vicinity" key inside SampleDataModel question object no. ' + i + ' not found.</br>';
                        lstQuestionDataModelFromFile.push(objQuestionFromFile);
                    }
                    objMainModelFromFile.SampleDataModel = lstQuestionDataModelFromFile;
                }

                mainModelFromFile.QuestionsModel = objMainModelFromFile;
                if (msg !== "") {
                    customException(msg);
                }
            }
            catch (err) {
                new PNotify({
                    title: 'Info!',
                    html: err,
                    type: 'info',
                    delay: 2000
                });

                //resetFileUplaodView();
                mainModelFromFile.ProjectId = "";
                return;
            }

            if (mainModelFromFile.ProjectId === "") {
                new PNotify({
                    title: 'Info!',
                    text: "Project with given Id is not exist.",
                    type: 'info',
                    delay: 2000
                });
            }
            else {
                $.ajax({
                    type: "POST",
                    url: generateHiveCall,
                    data: JSON.stringify(mainModelFromFile),
                    contentType: "application/json; charset=utf-8",
                    dataType: "json",
                    async: true,
                    beforeSend: function () {
                        showLoading();
                    },
                    complete: function () {
                        hideLoading()
                    },
                    success: function (data, textStatus, jqXHR) {
                        if (data === "LogOut") {
                            new PNotify({
                                title: 'Info!',
                                text: "Please Login to Submit! ",
                                type: 'info',
                                delay: 2000
                            }).then((result) => {
                                if (result.value) {
                                    var url = "/Home/Index";
                                    window.location.href = url;
                                }
                            });
                        }
                        else {
                            ajaxreturncode = jqXHR.status;
                            new PNotify({
                                type: 'success',
                                title: 'Asset Created Successfully!',
                                delay: 2000
                            });
                            resetFileUplaodView();
                            mainModelFromFile.ProjectId = "";
                        }
                    },
                    error: function (data, textStatus, jqXHR) {
                        ajaxreturncode = jqXHR.status;
                        new PNotify({
                            title: 'Error!',
                            text: textStatus,
                            type: 'error',
                            delay: 2000
                        });
                    }
                });
            }
        }, 2000);
    }
}

function showUploadButton() {
    showLoading();
    var reader = new FileReader();
    var file = $("#assetFile")[0].files[0];
    reader.readAsText(file);
    reader.onload = function (e) {
        try {
            var v = JSON.parse(e.target.result);
            document.getElementById('result').innerHTML = JSON.stringify(v, undefined, 2);
            $("#divResult").show();
            $("#uploadbutton").removeClass("hide");
            $('.footerDiv').css({ "position": "relative" });
            hideLoading();
        }
        catch (err) {
            hideLoading();
            new PNotify({
                title: 'Info!',
                text: err,
                type: 'info',
                delay: 2000
            });
            resetFileUplaodView();
        }
    }

}

function resetFileUplaodView() {
    $("#result").innerHTML = "";
    $("#divResult").hide();
    $("#assetFile").val("");
    $("#uploadbutton").addClass("hide");
    $('.footerDiv').css({ "position": "fixed" });
}

function showLoading() {
    $('#loader').removeClass("hide");
}

function hideLoading() {
    $('#loader').addClass("hide");
}



function LoadProjects() {

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
            for (var i = 0; i < data.aaData.length; i++) {
                options.push({ label: data.aaData[i].Name, title: data.aaData[i].Name, value: data.aaData[i].Id });
            }

            $('#drpProject').multiselect('dataprovider', options);
            $('#drpProject').multiselect('enable');


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


function optionAddRemoveAndCreditHtml(optCount) {
    var customHtml = '<input type="number" min=1 id="optionsCredits' + optCount + '" placeholder="Credits" class="optionsCredits form-control form-control-option-credit">' + optionAddRemoveButtonHtml(optCount);
    return customHtml;
}

function optionAddRemoveButtonHtml(rmAddButtonCount) {
    var customHtml = '<div class="classRemoveOption form-control form-control-option-deleteBtn" id="removeOption-' + rmAddButtonCount + '"><i class = "fa fa-remove"></i></div>' + '<div id="addButton" class="classAddOption form-control form-control-option-addBtn"><i class = "fa fa-plus"></i></div>';
    return customHtml;
}

function questionOptionHtml(questionCount) {
    var customHtml = '<input type="text" id="optQuestion' + questionCount + '" placeholder="Option text" class="optQuestion form-control form-control-option-text" >';
    return customHtml;
}

function refreshPage() {
    window.location.reload();
}