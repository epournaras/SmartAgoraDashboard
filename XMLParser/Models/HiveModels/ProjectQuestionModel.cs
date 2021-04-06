using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using XMLParser.Models;

namespace Models.HiveModels
{
    public class ProjectQuestionModel
    {
        public string ProjectId { get; set; }
        //public string TaskId { get; set; }
        public Questions QuestionsModel { get; set; }
    }
}
