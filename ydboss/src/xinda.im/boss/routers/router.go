package routers

import (
	"github.com/astaxie/beego"
	. "xinda.im/boss/controllers"
)

var (
	//登录
	loginCtl = &LoginController{}

	//销售
	accountCtl = &AccountController{}

	// 渠道
	AgentCtl = &AgentController{}

	// 客户
	CustomerCtl = &CustomerController{}

	// 商机
	BizCtl = &BizController{}
)

func init() {
	//登录
	beego.Router("/", loginCtl)
	beego.Router("/login.html", loginCtl)
	beego.Router("/login", loginCtl, "POST:Login")
	beego.Router("/loginOut", loginCtl, "GET:LoginOut")
	beego.Router("/YDLogin", loginCtl, "GET:YDLogin")

	// 获取Tag
	beego.Router("/get_all_tag", CustomerCtl, "POST:GetAllTag")

	// 关键字获取Tag
	beego.Router("/get_key_tag", CustomerCtl, "POST:GetTagByKey")

	//账号
	beego.Router("/employee.list", accountCtl, "POST:GetAccountList")
	beego.Router("/keyWord_employee", accountCtl, "POST:KeyWordAccount")
	beego.Router("/alert_employee", accountCtl, "post:AlertAccount")
	beego.Router("/del_employee", accountCtl, "post:DelAccount")
	beego.Router("/add_employee", accountCtl, "post:AddAccount")
	beego.Router("/sale_customer.html", accountCtl)
	beego.Router("/account_filter", accountCtl, "POST:FilterAccount")
	beego.Router("/update_account_psw", accountCtl, "POST:UpdateAccountPsw")
	beego.Router("/reset_account_psw", accountCtl, "POST:ResetAccountPsw")
	beego.Router("/employee_info.html", accountCtl, "GET:RenderEmployeeInfo")
	beego.Router("/get_custom_comments", accountCtl, "GET:GetCustomComments")
	beego.Router("/get_opport_comments", accountCtl, "GET:GetOpportComments")
	beego.Router("/agent_employee_sort", accountCtl, "POST:GetAccountBySort")

	// 渠道
	beego.Router("/agent.html", AgentCtl)
	beego.Router("/agent.list", AgentCtl, "POST:GetAgentList")
	beego.Router("/all_agent.list", AgentCtl, "POST:GetAgentList")
	beego.Router("/keyword_agent", AgentCtl, "post:KeyWordAgent")
	beego.Router("/alert_agent", AgentCtl, "post:AlertAgent")
	beego.Router("/del_agent", AgentCtl, "post:DelAgent")
	beego.Router("/add_agent", AgentCtl, "post:InsertAgent")
	beego.Router("/agent_employee.html", AgentCtl, "GET:AgentEmployee")
	beego.Router("/agent_customer.html", AgentCtl, "GET:AgentCustomer")
	beego.Router("/agent_filter", AgentCtl, "POST:FilterAgent")
	beego.Router("/agent_sort", AgentCtl, "POST:GetAgentBySort")

	// 客户
	beego.Router("/all_customer.html", CustomerCtl)
	beego.Router("/customer.list", CustomerCtl, "POST:GetAllCustomerList")
	beego.Router("/agent.customer.list", CustomerCtl, "POST:GetAgentCustomerList")
	beego.Router("/employee.customer.list", CustomerCtl, "POST:GetEmpCustomerList")

	beego.Router("/all_customer_sort", CustomerCtl, "POST:GetAllCustomerBySort")
	beego.Router("/sale_customer_sort", CustomerCtl, "POST:GetSaleCustomerBySort")
	beego.Router("/agent_customer_sort", CustomerCtl, "POST:GetAgentCustomerBySort")

	beego.Router("/show_customer_for_sale", CustomerCtl, "POST:GetCustomerListForSale")
	beego.Router("/allocation_customer", CustomerCtl, "POST:AllocationCustomer")
	beego.Router("/allocation_customer_agent", CustomerCtl, "POST:AllocationCustomerAgent")

	beego.Router("/distribution_customer", CustomerCtl, "post:DistributionCustomer")
	//修改客户信息
	beego.Router("/alert_customer", CustomerCtl, "post:AlertCustomer")

	beego.Router("/keyword_customer", CustomerCtl, "post:AdminKeyWordCustomer")
	beego.Router("/agent.keyword_customer", CustomerCtl, "post:AgentKeyWordCustomer")
	beego.Router("/employee.keyword_customer", CustomerCtl, "post:EmployeeKeyWordCustomer")

	beego.Router("/del_customer", CustomerCtl, "POST:DelCustomer")
	beego.Router("/add_customer", CustomerCtl, "POST:AddCustomer")

	beego.Router("/get_all_province", CustomerCtl, "POST:GetAllProvince")
	beego.Router("/get_city", CustomerCtl, "POST:GetCity")
	beego.Router("/get_province_by_key", CustomerCtl, "POST:GetProvinceByKey")
	beego.Router("/sale_customer_search.html", CustomerCtl, "GET:RenderCustomerSearch")
	beego.Router("/agent_customer_search.html", CustomerCtl, "GET:RenderCustomerSearch")
	beego.Router("/all_customer_search.html", CustomerCtl, "GET:RenderCustomerSearch")
	beego.Router("/sale_customer_type_filter", CustomerCtl, "POST:GetSaleCustomerTypeList")

	// 客户详细界面
	beego.Router("/get_super_customer_detail.html", CustomerCtl, "GET:GetSuperCustomerDetail")
	beego.Router("/get_agent_customer_detail.html", CustomerCtl, "GET:GetAgentCustomerDetail")
	beego.Router("/get_emp_customer_detail.html", CustomerCtl, "GET:GetSaleCustomerDetail")
	beego.Router("/getEmpCustomInfo", CustomerCtl, "GET:GetEmpCustomInfo")
	beego.Router("/getAgentCustomInfo", CustomerCtl, "GET:GetAgentCustomerInfo")
	beego.Router("/getSuperCustomInfo", CustomerCtl, "GET:GetSuperCustomerInfo")
	beego.Router("/get_sale_comments", CustomerCtl, "POST:GetSaleComments")

	// 客户备注
	beego.Router("/add_customer_comment", CustomerCtl, "POST:AddComment")
	beego.Router("/alert_customer_comment", CustomerCtl, "POST:AlertComment")

	// 客户筛选
	beego.Router("/customer_super_filter", CustomerCtl, "POST:FilterCustomerSuperList")
	beego.Router("/customer_agent_filter", CustomerCtl, "POST:FilterCustomerAgentList")
	beego.Router("/customer_emp_filter", CustomerCtl, "POST:FilterCustomerEmpList")

	beego.Router("/customer_excel_sale", CustomerCtl, "POST:ExcelCustomEmp")
	beego.Router("/download_excel_model", CustomerCtl, "GET:ExcelModelDownLoad")
	beego.Router("/customer_exists", CustomerCtl, "POST:CustomerIsExists")

	// 商机
	beego.Router("/opportunities.html", BizCtl)
	beego.Router("/agent_biz.html", BizCtl, "GET:RenderAgentBiz")
	beego.Router("/opport.list", BizCtl, "POST:GetBizunitiesList")
	beego.Router("/agent_opport_list", BizCtl, "POST:AgentShowBizs")
	beego.Router("/agent_keyword_opport", BizCtl, "POST:GetAgentOpportKeyword")
	beego.Router("/keyword_opport", BizCtl, "post:KeyWordBizList")
	beego.Router("/alert_opport", BizCtl, "post:AlertBiz")
	beego.Router("/del_opport", BizCtl, "post:DelBiz")
	beego.Router("/add_opport", BizCtl, "post:InsertBiz")
	beego.Router("/opport_sort", BizCtl, "POST:GetOpportBySort")
	beego.Router("/agent_opport_sort", BizCtl, "POST:GetAgentOpportBySort")
	beego.Router("/sale_opport_search.html", BizCtl, "GET:RenderOpportSearch")
	beego.Router("/agent_opport_search.html", BizCtl, "GET:RenderOpportSearch")

	// 商机筛选
	beego.Router("/opport_super_filter", BizCtl, "POST:FilterBizEmpList")
	beego.Router("/opport_agent_filter", BizCtl, "POST:FilterBizAgentList")
	beego.Router("/opport_emp_filter", BizCtl, "POST:FilterBizEmpList")

	// 商机详细页
	beego.Router("/get_emp_biz_detail.html", BizCtl, "GET:GetEmpBizDetail")
	beego.Router("/get_agent_biz_detail.html", BizCtl, "GET:GetAgentBizDetail")
	beego.Router("/getEmpBizInfo", BizCtl, "GET:EmpBizInfo")
	beego.Router("/getBizComments", BizCtl, "POST:GetBizComments")
	beego.Router("/add_opport_comment", BizCtl, "POST:InsertCommentsFromC")
	beego.Router("/alert_biz_comment", BizCtl, "POST:AlertCommentsFromC")

}
