namespace HiveServer.Models
{
    public class Record
    {
        public string[] sensors { get; set; }
        public string start { get; set; }
        public string end { get; set; }
        public int step { get; set; }
    }
}