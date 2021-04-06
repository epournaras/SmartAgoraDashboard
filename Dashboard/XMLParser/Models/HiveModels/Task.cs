namespace HiveServer.Models
{
    public class Task
    {
        public string Id { get; set; }
        public string Name { get; set; }
        public string Description { get; set; }
        public string CurrentState { get; set; }
        public AssignmentCriteria AssignmentCriteria { get; set; }
        public CompletionCriteria CompletionCriteria { get; set; }
    }
    
}