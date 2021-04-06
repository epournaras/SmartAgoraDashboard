using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Web.Mvc;
namespace XMLParser.Filter
{
    public class CustomActionFilter: ActionFilterAttribute,IActionFilter
    {
        
        public override void OnActionExecuting(ActionExecutingContext filterContext)
        {
            var Session = HttpContext.Current.Session;

            if (System.Web.HttpContext.Current.Session["Email"] == null)    // User Not Logged In. Redirect or show an error
            {
                var control = filterContext.Controller as Controller;
                var controller = filterContext.RouteData.Values["controller"].ToString();
                var action = filterContext.RouteData.Values["action"].ToString();

                if (control != null)
                {
                    if (filterContext.HttpContext.Request.Url.ToString().Contains("Home") || filterContext.HttpContext.Request.Url.ToString().Contains("Hive"))
                    {
                        if (filterContext.HttpContext.Request.IsAjaxRequest())
                        {
                            var result = new JsonResult();
                            result.Data = "LogOut";
                            result.JsonRequestBehavior = JsonRequestBehavior.AllowGet;
                            filterContext.Result = result;
                            return;
                        }
                        control.HttpContext.Response.Clear();
                        control.HttpContext.Response.Redirect("/Home/Index");
                        control.HttpContext.Response.Close();
                    }
                }
            }

            
        }
        

        public override void OnActionExecuted(ActionExecutedContext filterContext)
        {
        }

        public override void OnResultExecuting(ResultExecutingContext filterContext)
        {
        }

        public override void OnResultExecuted(ResultExecutedContext filterContext)
        {
        }

    }
}