namespace HiveServer.Models
{
    public class Assignment
    {
        public string Id { get; set; }
        public string User { get; set; }
        public string Project { get; set; }
        public string Task { get; set; }
        public Asset Asset { get; set; }
        public string State { get; set; }
        public SubmittedData SubmittedData { get; set; }
        public SubmittedAnswerData SubmittedAnswerData { get; set; }
    }
}
