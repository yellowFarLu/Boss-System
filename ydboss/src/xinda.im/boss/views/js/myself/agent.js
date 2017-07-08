"user struct";
(function (Dom) {

    /*声明*/
    var GlobalFunc = new Global();
    var AgentCtl = new Agent();

    var Tpl_page = _.template(GlobalFunc.pageTpl),
        Tpl_agent_list = _.template($('#TplAgentList').html().replace(/\&lt\;/ig, '\<')),
        Tpl_Province_list = _.template($('#TplProvinceList').html().replace(/\&lt\;/ig, '\<')),
        Tpl_City_list = _.template($('#TplCityList').html().replace(/\&lt\;/ig, '\<')),
        $addAgentBtn = $('#addAgentBtn'),
        $searchAgentBtn = $('#searchAgentBtn'),
        $ModalChangeAgent = $('#ModalChangeAgent'),
        $changeAgentBtn = $('#changeAgentBtn'),
        $delAgentBtn = $('#DelAgentBtn'),
        $agent_page = $('#agent_page'),
        $addAgentForm = $('#addAgentForm'),
        $updateAgentForm = $('#updateAgentForm'),
        $ModalAddAgent = $('#ModalAddAgent'),
        $ModalChangeAgent = $('#ModalChangeAgent'),
        $delAgentNow = $('#delAgentNow'),
        $ModalDelAgent = $('#ModalDelAgent'),
        $addAgentNow = $('#addAgentNow')

    var That;

    GlobalFunc.Focus('#addAgentNow', '#addAgentForm');
    GlobalFunc.Focus('.changeAgentBtn', '#updateAgentForm');

    function Agent() {
        this.agent = {};
        this.page_info = null;
        this.keyOn = null;
        this.key = null;
        this.filter = null;
        this.sort = null;
        That = this;

        this.showAgentList = function (data) {
            var h_table = Tpl_agent_list(data);
            var h_page = Tpl_page(data.page_info);
            $('#agentListBox').html(h_table);
            $('#agent_page').html(h_page);
        }
        this.changeAgent = function () {
            this.agent.AgentName = Dom.find('#_AgentName').val();
            this.agent.Contacts = Dom.find('#_Contacts').val();
            this.agent.ContactsMail = Dom.find('#_ContactsMail').val();
            this.agent.ContactsPhone = Dom.find('#_ContactsPhone').val();
            this.agent.Note = Dom.find('#_Note').val();
            this.agent.AgentId = Dom.find('#_AgentId').val();
        }
        this.Flip = function (page) {
            switch (this.keyOn) {
            case 1:
                this.searchAgent(this.key, page);
                break;
            case 2:
                var o = this.sort;
                o.page = page;
                this.sortByKey(o);
                break;
            default:
                this.getAgentList(page)
            }
        }
    }

    Agent.prototype.getAgentList = function (page) {
        var p = {
            page: page,
            all: false
        }
        $.post('/agent.list', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 0;
                That.showAgentList(data);
            }
        })
    }

    Agent.prototype.getCity = function (id, type, derfer) {
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

    Agent.prototype.addAgent = function (agent) {
        $.post('/add_agent', agent, function (data) {
            if (!!data && data.code == 0) {
                That.getAgentList(1);
                $ModalAddAgent.modal('hide');
                $addAgentForm[0].reset();
                GlobalFunc.Alert.success('新增成功 \^o^/');
            } else if (!!data && data.code == 3) {
                GlobalFunc.Alert.info('改账号已存在，请重新输入新账号 X﹏X');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Agent.prototype.searchAgent = function (key, page) {
        var p = {
            keyWord: key,
            page: page || 1
        }
        $.post('/keyword_agent', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.key = key;
                That.keyOn = 1;
                That.showAgentList(data)
            }
        })
    }

    Agent.prototype.deleteAgent = function (ids) {
        var p = {
            AgentIds: ids
        }
        $.post('/del_agent', p, function (data) {
            if (!!data && data.code == 0) {
                That.getAgentList(1);
                GlobalFunc.Alert.success('删除成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Agent.prototype.updateAgent = function (agent) {
        $.post('/alert_agent', agent, function (data) {
            if (!!data && data.code == 0) {
                That.getAgentList(1);
                $ModalChangeAgent.modal('hide');
                GlobalFunc.Alert.success('更新成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

//    Agent.prototype.agentFilter = function (o) {
//        $.post('/agent_filter', o, function (data) {
//            if (!!data && data.code == 0) {
//                That.page_info = data.page_info;
//                That.keyOn = 2;
//                That.filter = o;
//                That.showAgentList(data);
//            }
//        })
//    }
    
    Agent.prototype.sortByKey = function (p) {
        $.post('/agent_sort', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 2;
                That.sort = p;
                That.showAgentList(data);
            }
        })
    }


    AgentCtl.getAgentList(1);

    /*时间绑定*/
    Dom.find('.modal').modal({
        backdrop: 'static',
        keyboard: false,
        show: false
    });

    $addAgentBtn.on('click', function () {
        if (GlobalFunc.CheckForm($addAgentForm)) return;
        var agent = new Object();
        agent.AgentName = Dom.find('#AgentName').val();
        agent.Contacts = Dom.find('#Contacts').val();
        agent.ContactsMail = Dom.find('#ContactsMail').val();
        agent.ContactsPhone = Dom.find('#ContactsPhone').val();
        agent.Note = Dom.find('#Note').val();

        agent.AccountId = Dom.find('#AccountId').val();
        agent.Name = Dom.find('#Contacts').val();
        agent.Password = Dom.find('#Password').val();
        agent.Sex = Dom.find('#Sex').find('option:checked').val();

        AgentCtl.addAgent(agent);
    });

    $searchAgentBtn.on('click', function () {
        AgentCtl.keyOn = 1;
        var key = Dom.find('#keyWord').val();
        if (!$.trim(key)) {
            AgentCtl.getAgentList(1);
            return
        }
        AgentCtl.searchAgent(key)
    });

    $delAgentBtn.on('click', function () {
        var chk_value = [];
        $('input[name="DelAgentCheck"]:checked').each(function () {
            chk_value.push($(this).val());
        });
        AgentCtl.deleteAgent(chk_value)
    });

    Dom.on('click', '.changeAgentBtn', function () {
        $ModalChangeAgent.modal('show');
        var parent_dom = $(this).parents('tr');
        var agent = new Object();
        agent.AgentName = parent_dom.find('.AgentName').text();
        agent.Contacts = parent_dom.find('.Contacts').text();
        agent.ContactsMail = parent_dom.find('.ContactsMail').text();
        agent.ContactsPhone = parent_dom.find('.ContactsPhone').text();
        agent.Note = parent_dom.find('.Note').text();
        agent.AgentId = parent_dom.find('.AgentId').data('id');

        Dom.find('#_AgentName').val(agent.AgentName);
        Dom.find('#_Contacts').val(agent.Contacts);
        Dom.find('#_ContactsMail').val(agent.ContactsMail);
        Dom.find('#_ContactsPhone').val(agent.ContactsPhone);
        Dom.find('#_Note').val(agent.Note);
        Dom.find('#_AgentId').val(agent.AgentId);

    });

    $changeAgentBtn.on('click', function () {
        if (GlobalFunc.CheckForm($updateAgentForm)) return;
        AgentCtl.changeAgent();
        AgentCtl.updateAgent(AgentCtl.agent);
    });

    $delAgentNow.on('click', function () {
        !GlobalFunc.CheckChoose('DelAgentCheck') ? $ModalDelAgent.modal('show') : GlobalFunc.Alert.info('请勾选你需要删除的内容！')
    });

    $addAgentNow.click(function () {
        $ModalAddAgent.modal('show');
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
        AgentCtl.sortByKey(sort);
    });


    /*排序*/
//    Dom.on('click', '#timexSortBtn', function () {
//        var dom = $(this).parents('.panel-body');
//        var sortIndex = dom.find('.orderBy').find('input[type=radio]:checked').data('index');
//
//        var o = {
//            sortIndex: sortIndex,
//            sortKey: 'timex',
//            timex_begin: dom.find('.startValue').val(),
//            timex_end: dom.find('.endValue').val()
//        }
//
//        AgentCtl.agentFilter(o);
//    })

    /*翻页*/
    $agent_page.on('click', '.trun-left', function () {
        var page = AgentCtl.page_info.page - 1 || 1;
        AgentCtl.Flip(page)
    })
    $agent_page.on('click', '.turn-right', function () {
        var page = AgentCtl.page_info.page = AgentCtl.page_info.page + 1 > AgentCtl.page_info.page_total ? AgentCtl.page_info.page_total : AgentCtl.page_info.page + 1;
        AgentCtl.Flip(page)
    })
    $agent_page.on('click', '.branch-btn', function () {
        var page = Number($(this).text());
        AgentCtl.Flip(page)
    })


})($(document.body))