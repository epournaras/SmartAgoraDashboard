using System;
using System.Collections.Generic;
namespace HiveServer.Models
{

public class Options
{
    public string Credits { get; set; }
    public string Name { get; set; }
    public object NextQuestion { get; set; }
    public int id { get; set; }
}

public class Sensors
{
    public string Name { get; set; }
    public int id { get; set; }
}

public class SampleDataModels
{
    public object Combination { get; set; }
    public string Frequency { get; set; }
    public string Latitude { get; set; }
    public string Longitude { get; set; }
    public string Mandatory { get; set; }
    public List<Options> Option { get; set; }
    public string Question { get; set; }
    public List<Sensors> Sensor { get; set; }
    public string Sequence { get; set; }
    public string Time { get; set; }
    public string Type { get; set; }
    public string Vicinity { get; set; }
    public string Visibility { get; set; }
    public int id { get; set; }
}

public class StartAndDestinationModels
{
    public string DefaultCredit { get; set; }
    public object DestinationLatitude { get; set; }
    public object DestinationLongitude { get; set; }
    public string Mode { get; set; }
    public object StartLatitude { get; set; }
    public object StartLongitude { get; set; }
}

public class Records
{
    public List<SampleDataModels> SampleDataModel { get; set; }
    public StartAndDestinationModels StartAndDestinationModel { get; set; }
}

public class Metadatas
{
    public Record record { get; set; }
}

public class SubmittedDatas
{
    public object test5 { get; set; }
}

public class Counts
{
    public int Assignments { get; set; }
    public int finished { get; set; }
    public int skipped { get; set; }
    public int unfinished { get; set; }
}

public class Assets
{
    public string Id { get; set; }
    public string Project { get; set; }
    public string Url { get; set; }
    public string Name { get; set; }
    public Metadatas Metadata { get; set; }
    public SubmittedDatas SubmittedData { get; set; }
    public bool Favorited { get; set; }
    public bool Verified { get; set; }
    public Counts Countss { get; set; }
}

public class Assignments
{
    public string Id { get; set; }
    public string User { get; set; }
    public string Project { get; set; }
    public string Task { get; set; }
    public Assets Asset { get; set; }
    public string State { get; set; }
    public object SubmittedData { get; set; }
}

}
