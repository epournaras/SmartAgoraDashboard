namespace HiveServer.Models
{
    public class Project
    {
        public string Id { get; set; }
        public string Name { get; set; }
        public string Description { get; set; }
        public int AssetCount { get; set; }
        public int TaskCount { get; set; }
        public int UserCount { get; set; }
        public MetaProperty[] MetaProperties { get; set; }
        public AssignmentCount AssignmentCount { get; set; }
    }
    public class AssignmentCount
    {
        public int Total { get; set; }
        public int Finished { get; set; }
    }
}