"user struct";
(function (Dom) {

    /*声明*/
    var GlobalFunc = new Global();
    var CustomerCtl = new Customer();
    var Tpl_page = _.template(GlobalFunc.pageTpl),
        Tpl_Customer_list = _.template($('#TplCustomerList').html().replace(/\&lt\;/ig, '\<')),
        Tpl_Agent_list = _.template($('#TplAgentList').html().replace(/\&lt\;/ig, '\<')),
        Tpl_Province_list = _.template($('#TplProvinceList').html().replace(/\&lt\;/ig, '\<')),
        Tpl_City_list = _.template($('#TplCityList').html().replace(/\&lt\;/ig, '\<')),
        $addCustomerBtn = $('#addCustomerBtn'),
        $searchCustomerBtn = $('#searchCustomerBtn'),
        $updateCustomerBtn = $('#updateCustomerBtn'),
        $ModalChangeCustomer = $('#ModalChangeCustomer'),
        $delCustomerBtn = $('#DelCustomerBtn'),
        $Customer_page = $('#Customer_page'),
        $AgentList = $('#AgentList'),
        $AllocationBtn = $('#AllocationBtn'),
        $addCustomerForm = $('#addCustomerForm'),
        $ModalAddCustomer = $('#ModalAddCustomer'),
        $updateCustomerForm = $('#updateCustomerForm'),
        $ModalChangeCustomer = $('#ModalChangeCustomer'),
        $delCustomerNow = $('#delCustomerNow'),
        $ModalDelCustomer = $('#ModalDelCustomer'),
        $AllocationCustomerNow = $('#AllocationCustomerNow'),
        $ModalAllocation = $('#ModalAllocation'),
        $addCustomerNow = $('#addCustomerNow')
        

    var That;
    
    GlobalFunc.Focus('#addCustomerNow','#addCustomerForm');
    GlobalFunc.Focus('.changeCustomerBtn','#updateCustomerForm');

    function Customer() {
        this.Customer = {};
        this.page_info = null;
        this.keyOn = null;
        this.key = null;
        this.filter = null;
        this.filter_info = {};
        this.sort = null;
        That = this;

        this.showCustomerList = function (data) {            
            var h_table = Tpl_Customer_list(data);
            var h_page = Tpl_page(data.page_info);
            $('#CustomerListBox').html(h_table);
            $Customer_page.html(h_page);
            $('.popovers').popover({
                html: true
            });
        }
        this.changeCustomer = function () {
            this.Customer.rtx_num = Dom.find('#_Buin').val();
            this.Customer.entName = Dom.find('#_EntName').val();
            this.Customer.contacts = Dom.find('#_Contacts').val();
            this.Customer.mail = Dom.find('#_ContactsMail').val();
            this.Customer.phone = Dom.find('#_ContactsPhone').val();
            this.Customer.mobile = Dom.find('#_Mobile').val();
            this.Customer.qq = Dom.find('#_QQ').val();
            this.Customer.province = Dom.find('#_province').find('option:checked').val();
            this.Customer.city = Dom.find('#_city').find('option:checked').val();
            // this.Customer.follow_time = Dom.find('#_follow_time').val();
            this.Customer.remarks = Dom.find('#_Remarks').val();
            var tags = [];
            Dom.find('#_updateTags').find('span').each(function () {
                tags.push($(this).data('id'));
            });
            this.Customer.tags = tags;
        }
        this.Flip = function (page) {
            switch (this.keyOn) {
            case 1:
                this.searchCustomer(this.key, page);
                break;
            case 2:
                var o = this.sort;
                o.page = page;
                this.sortByKey(o);
                break;
            default:
                this.getCustomerList(page)
            }
        }
    }

    Customer.prototype.getCustomerList = function (page) {
        var p = {
            page: page
        }
        $.post('/agent.customer.list', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 0;
                That.showCustomerList(data);
            }
        })
    }     
    
    Customer.prototype.getAllProvince = function () {
        $.post('/get_all_province', function (data) {
            if (!!data && data.code == 0) {
                var h = Tpl_Province_list(data);
                $('#province,#_province').html(h);
            }
        })
    }

    Customer.prototype.getCity = function (id, type, derfer) {
        var p = {
            province_id: id
        }
        $.post('/get_city', p, function (data) {
            if (!!data && data.code == 0) {
                var h = Tpl_City_list(data);
                type == 0 ? $('#city').html(h) : $('#_city').html(h);
                derfer && derfer.resolve();
            }
        })
        var rsp = derfer ? derfer.promise() : 0;
        return rsp;
    }   
    
    Customer.prototype.getAllAgentList = function (page) {
        var p = {
            page: page,
            all: true
        }
        $.post('/employee.list', p, function (data) {
            if (!!data && data.code == 0) {
                var h = Tpl_Agent_list(data)
                $AgentList.html(h);
            }
        })
    }

    Customer.prototype.addCustomer = function (Customer) {
        Customer.url_type = 1;
        $.post('/add_customer', Customer, function (data) {
            if (!!data && data.code == 0) {
                That.getCustomerList(1);
                
                $addCustomerForm[0].reset();
                $ModalAddCustomer.modal('hide');
                
                GlobalFunc.Alert.success('新增成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Customer.prototype.searchCustomer = function (key, page) {
        var p = {
            keyWord: key,
            page: page || 1
        }
        $.post('/agent.keyword_customer', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.key = key;
                That.keyOn = 1;
                That.showCustomerList(data);
            }
        })
    }

    Customer.prototype.deleteCustomer = function (ids) {
        var p = {
            CustomerIds: ids
        }
        $.post('/del_customer', p, function (data) {
            if (!!data && data.code == 0) {
                That.getCustomerList(1);
                GlobalFunc.Alert.success('删除成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Customer.prototype.updateCustomer = function (Customer) {
        $.post('/alert_customer', Customer, function (data) {
            if (!!data && data.code == 0) {
                That.getCustomerList(1);
                GlobalFunc.Alert.success('更新成功 \^o^/');
                $updateCustomerForm[0].reset();
                $ModalChangeCustomer.modal('hide');

            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Customer.prototype.customFilter = function (o) {
        $.post('/customer_agent_filter', o, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 2;
                That.filter = o;
                That.showCustomerList(data);
            }
        })
    }
    
    Customer.prototype.allocationCustomer = function (ids, accId) {
        var p = {
            CustomerIds: ids,
            EmployeeId : accId
        }
        $.post('/allocation_customer', p, function (data) {
            if (!!data && data.code == 0) {
                That.getCustomerList(1);
                GlobalFunc.Alert.success('分配成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }
    
    
    Customer.prototype.sortByKey = function (p) {
        $.post('/agent_customer_sort', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 2;
                That.sort = p;
                That.showCustomerList(data);
            }
        })
    }
    
/* ---------------------------------------------------------------------------------------------------------------------------------------------------- */

    CustomerCtl.getCustomerList(1);
    CustomerCtl.getAllAgentList(1);
    CustomerCtl.getAllProvince();

    /*时间绑定*/
    Dom.find('.modal').modal({
        backdrop: 'static',
        keyboard: false,
        show: false
    });

    $addCustomerBtn.on('click', function () {
        if (GlobalFunc.CheckForm($addCustomerForm) || GlobalFunc.CheckArea('#province', '#city')) {
            return
        };

        var Customer = new Object();
        Customer.buin = Dom.find('#Buin').val();
        Customer.entName = Dom.find('#EntName').val();
        Customer.contacts = Dom.find('#Contacts').val();
        Customer.mail = Dom.find('#ContactsMail').val();
        Customer.phone = Dom.find('#ContactsPhone').val();
        Customer.mobile = Dom.find('#Mobile').val();
        Customer.qq = Dom.find('#QQ').val();
        Customer.province = Dom.find('#province').find('option:checked').val();
        Customer.city = Dom.find('#city').find('option:checked').val();
        Customer.remarks = Dom.find('#Remarks').val();
        
        var tags = [];
        Dom.find('#tagAll').find('.tag.on').each(function () {
            tags.push($(this).find('span').data('id'));    
        });
        Customer.tags = tags;

        CustomerCtl.addCustomer(Customer);
    });

    $searchCustomerBtn.on('click', function () {
        CustomerCtl.keyOn = true;
        var key = Dom.find('#keyWord').val();
        if (!$.trim(key)) {
            CustomerCtl.getCustomerList(1);
            return
        }
        CustomerCtl.searchCustomer(key)
    })

    $delCustomerBtn.on('click', function () {
        var chk_value = [];    
        $('input[name="DelCustomerCheck"]:checked').each(function () {
            chk_value.push($(this).val());   
        });
        CustomerCtl.deleteCustomer(chk_value)
    });
    
    $AllocationBtn.on('click', function () {
        var chk_value = [];
        $('input[name="DelCustomerCheck"]:checked').each(function () {
            chk_value.push($(this).val());   
        });
        var id = $AgentList.find('option:checked').val();
        CustomerCtl.allocationCustomer(chk_value, id)
    })
    
    Dom.on('click', '.changeCustomerBtn', function () {     
        $ModalChangeCustomer.modal('show');
        var parent_dom = $(this).parents('tr');
        var Customer = new Object();
        var der = $.Deferred();
        var d;
        Customer.buin = parent_dom.find('.buin').text() == 0 ? "" : parent_dom.find('.buin').text();
        Customer.entName = parent_dom.find('.entName').text();
        Customer.contacts = parent_dom.find('.contacts').text();
        Customer.mail = parent_dom.find('.email').text();
        Customer.phone = parent_dom.find('.phone').text();
        Customer.mobile = parent_dom.find('.mobile').text();
        Customer.province = parent_dom.find('.province').text();
        Customer.city = parent_dom.find('.city').text();
        Customer.qq = parent_dom.find('.qq').text();
        Customer.remarks = parent_dom.find('.remarks').text();
        Customer.customerId = parent_dom.find('.customerId').data('id');
        Customer.follow_time = parent_dom.find('.follow_time').text();
        Customer.tags_html = parent_dom.find('.tag').html();
        
        $('#_tagAll .tag').removeClass('on');
        parent_dom.find('.tag').children('span').each(function () {
            var _class = ".tag_" + $(this).data('id');
            $('#_tagAll').find(_class).parent().addClass('on');
        });

        Dom.find('#_Buin').val(Customer.buin);
        Dom.find('#_EntName').val(Customer.entName);
        Dom.find('#_Contacts').val(Customer.contacts);
        Dom.find('#_ContactsMail').val(Customer.mail);
        Dom.find('#_ContactsPhone').val(Customer.phone);
        Dom.find('#_Mobile').val(Customer.mobile);
        Dom.find('#_QQ').val(Customer.qq);
        Dom.find('#_follow_time').val(Customer.follow_time);
        Dom.find('#_updateTags').html(Customer.tags_html);

        Dom.find('#_province').find('option').each(function () {
            if ($(this).text() == Customer.province) {
                $(this)[0].selected = true;
                var id = $(this).val();
                d = CustomerCtl.getCity(+id, 1, der);
            }
        });
        
        d.then(function () {
            Dom.find('#_city').find('option').each(function () {
                $(this).text() == Customer.city && ($(this)[0].selected = true);
            });
        });
        
        Dom.find('#_Remarks').val(Customer.remarks);
        CustomerCtl.Customer.customerId = Customer.customerId;
    });

    $updateCustomerBtn.on('click', function () {
         if (GlobalFunc.CheckForm($updateCustomerForm) || GlobalFunc.CheckArea('#_province', '#_city')) {
            return
        };
        
        CustomerCtl.changeCustomer();
        CustomerCtl.updateCustomer(CustomerCtl.Customer); 
    });
    
    $delCustomerNow.on('click', function () {
        !GlobalFunc.CheckChoose('DelCustomerCheck') ? $ModalDelCustomer.modal('show') : GlobalFunc.Alert.info('请勾选你需要删除的客户！');
    });
    
    $AllocationCustomerNow.on('click', function () {
        !GlobalFunc.CheckChoose('DelCustomerCheck') ? $ModalAllocation.modal('show') : GlobalFunc.Alert.info('请勾选你需要分配的客户！');
    });

    Dom.find('#province').change(function () {
        var id = $(this).find('option:checked').val();
        if (!id) return
        CustomerCtl.getCity(+id, 0)
    });

    Dom.find('#_province').change(function () {
        var id = $(this).find('option:checked').val();
        if (!id) return
        CustomerCtl.getCity(+id, 1)
    });
    
    $addCustomerNow.click(function () {
        $ModalAddCustomer.modal('show');
    });
    
    /**
     * Sort Button [Onclick]
     */
    Dom.on('click', '.orderByBtn', function () {
        var $this = $(this);
        if ($this.hasClass('on')) {
            var DOM_I = $this.find('i.fa.on');
            DOM_I.removeClass('on').siblings('i').addClass('on');
        } else {
            $this.find('i.fa:eq(1)').addClass('on');
        }
        $this.siblings('.orderByBtn').removeClass('on').find('i.fa').removeClass('on');
        $this.addClass('on');
        
        var sort = {
            sortKey: $this.data('key'),
            sortIndex: $this.find('i.fa.on').index(),
            page: 1
        };
        CustomerCtl.sortByKey(sort);
    });

    /*排序*/
//    Dom.on('click', '#timexSortBtn, #followTimeSortBtn', function () {
//        var dom = $(this).parents('.panel-body').first();
//        var key = $(this).data('key');
//        var sortIndex = dom.find('.orderBy').find('input[type=radio]:checked').data('index');
//
//        CustomerCtl.filter_info.sortIndex = sortIndex;
//        CustomerCtl.filter_info.sortKey = key;
//        if (key == 'timex') {
//            CustomerCtl.filter_info.timex_begin = dom.find('.startValue').val();
//            CustomerCtl.filter_info.timex_end = dom.find('.endValue').val();
//        } else {
//            CustomerCtl.filter_info.last_follow_time_begin = dom.find('.startValue').val();
//            CustomerCtl.filter_info.last_follow_time_end = dom.find('.endValue').val();
//        }
//
//        CustomerCtl.customFilter(CustomerCtl.filter_info);
//    });
//
//    Dom.on('click', '#tagSortBtn', function () {
//        var dom = $(this).parents('.panel-body').first();
//        var tags = [];
//        dom.find('#allTagsBox').children('.tag.on').each(function () {
//            tags.push($(this).data('id'));
//        });
//        CustomerCtl.filter_info.tags = tags;
//        CustomerCtl.customFilter(CustomerCtl.filter_info);
//    });
//
//    Dom.on('click', '#areaSortBtn', function(){
//        var dom = $(this).parents('.panel-body').first();
//        var privinces = [];
//        dom.find('#allAreaBox').children('.tag.on').each(function () {
//            privinces.push($(this).data('id'));
//        });
//        CustomerCtl.filter_info.provinces = privinces;
//        CustomerCtl.customFilter(CustomerCtl.filter_info);
//    });


    /*翻页*/
    $Customer_page.on('click', '.trun-left', function () {
        var page = CustomerCtl.page_info.page - 1 || 1;
        CustomerCtl.Flip(page);
    });
    $Customer_page.on('click', '.turn-right', function () {
        var page = CustomerCtl.page_info.page = CustomerCtl.page_info.page + 1 > CustomerCtl.page_info.page_total ? CustomerCtl.page_info.page_total : CustomerCtl.page_info.page + 1;
        CustomerCtl.Flip(page);
    });
    $Customer_page.on('click', '.branch-btn', function () {
        var page = Number($(this).text());
        CustomerCtl.Flip(page);
    });


})($(document.body))