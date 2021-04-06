using System;
using System.Collections.Generic;


namespace HiveServer.Models
{
    public class Asset
    {
        public string Id { get; set; }
        public string Name { get; set; }
        public string Url { get; set; }
        public Metadata Metadata { get; set; }

        //public static implicit operator List<object>(Asset v)
        //{
        //    throw new NotImplementedException();
        //}
        //public object Metadata { get; set; }
        //public submittedData SubmittedData { get; set; }

        //public bool Favorited { get; set; }
        //public bool Verified { get; set; }
        //public Count Counts { get; set; }

    }
}