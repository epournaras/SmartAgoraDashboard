namespace HiveServer.Models
{
    public class CompletionCriteria
    {
        public int Total { get; set; }
        public int Matching { get; set; }

    }
    public class SocialExperimentXML
    {
        public string Data { get; set; }
    }

    public class Count
    {
        public string Assignments { get; set; }
        public int finished { get; set; }
        public int skipped { get; set; }
        public int unfinished { get; set; }
    }
 

    public class submittedData
    {
        public string Data { get; set; }
    }
}