"user struct";
(function (Dom) {

    /*声明*/
    var GlobalFunc = new Global();
    var OpportCtl = new Opport();
    var Tpl_page = _.template(GlobalFunc.pageTpl),
        Tpl_Opport_list = _.template($('#TplOpportList').html()),
        TplCustomerList = _.template($('#TplCustomerList').html()),
        $searchOpportBtn = $('#searchOpportBtn'),
        $ModalChangeOpport = $('#ModalChangeOpport'),
        $changeOpportBtn = $('#changeOpportBtn'),
        $delOpportBtn = $('#DelOpportBtn'),
        $Opport_page = $('#opport_page'),
        $updateOpportForm = $('#updateOpportForm'),
        $ModalChangeOpport = $('#ModalChangeOpport'),
        $ModalDelOpport = $('#ModalDelOpport'),
        $closeOpportBtn = $('#closeOpportBtn')

    var That;
    
    GlobalFunc.Focus('.alertOpportBtn','#updateOpportForm');

    function Opport() {
        this.Opport = {};
        this.page_info = null;
        this.keyOn = null;
        this.key = null;
        this.list_length = 0;
        this.filter_info = {};
        this.sort = null;
        That = this;

        this.showOpportList = function (data) {
            _.map(data.list, function (v) {
                v.status = GlobalFunc.OpportStage[v.status];
            });
            var h_table = Tpl_Opport_list(data);
            var h_page = Tpl_page(data.page_info);
            $('#opportListBox').html(h_table);
            $('#opport_page').html(h_page);
        }
        this.changeOpport = function () {
            this.Opport.title = Dom.find('#_OpportTitle').val();
            this.Opport.content = Dom.find('#_Content').val();
            this.Opport.status = Dom.find('#_Stage').find('option:checked').val();
            this.Opport.type = Dom.find('#_Type').find('option:checked').val();
            this.Opport.quota = Dom.find('#_Quota').val();
            this.Opport.opportType = Dom.find('#_OpportType').find('option:checked').val();
            this.Opport.real_time = Dom.find('#_RealTime').val();
            
            var tags = [];
            Dom.find('#_updateTags').find('span').each(function () {
                tags.push($(this).data('id'));
            });

            this.Opport.tags = tags;
        }
        this.Flip = function (page) {
            switch (this.keyOn) {
            case 1:
                this.searchOpport(this.key, page);
                break;
            case 2:
                var o = this.sort;
                o.page = page;
                this.sortByKey(o);
                break;
            default:
                this.getOpportList(page)
            }
        }
    }

    Opport.prototype.getOpportList = function (page) {
        var p = {
            page: page,
            all: false
        }
        $.post('/agent_opport_list', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 0;
                That.showOpportList(data);
            }
        })
    }

    Opport.prototype.searchOpport = function (key, page) {
        var p = {
            keyWord: key,
            page: page || 1
        }
        $.post('/agent_keyword_opport', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.key = key;
                That.keyOn = 1;
                That.showOpportList(data);
            }
        })
    }

    Opport.prototype.deleteOpport = function (ids) {
        var p = {
            OpportIds: ids
        }
        $.post('/del_opport', p, function (data) {
            if (!!data && data.code == 0) {
                That.getOpportList(1);
                GlobalFunc.Alert.success('关闭成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Opport.prototype.updateOpport = function (Opport) {
        $.post('/alert_opport', Opport, function (data) {
            if (!!data && data.code == 0) {
                That.getOpportList(1);
                $ModalChangeOpport.modal('hide');
                $updateOpportForm[0].reset();
                GlobalFunc.Alert.success('更新成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }
    
    
//    Opport.prototype.bizFilter = function(o) {
//        $.post('/biz_filter', o, function(data) {
//            if (!!data && data.code == 0) {
//                That.page_info = data.page_info;
//                That.keyOn = 2;
//                That.filter = o;
//                That.showOpportList(data);
//            }
//        })
//    }
    
    
    Opport.prototype.sortByKey = function (p) {
        $.post('/agent_opport_sort', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 2;
                That.sort = p;
                That.showOpportList(data);
            }
        })
    }


    OpportCtl.getOpportList(1);

    /*时间绑定*/
    Dom.find('.modal').modal({
        backdrop: 'static',
        keyboard: false,
        show: false
    });

    $searchOpportBtn.on('click', function () {
        OpportCtl.keyOn = true;
        var key = Dom.find('#keyWord').val();
        if (!$.trim(key)) {
            OpportCtl.getOpportList(1);
            return
        }
        OpportCtl.searchOpport(key)
    })

    $delOpportBtn.on('click', function () {
        var chk_value = []; //定义一个数组      
        $('input[name="DelOpportCheck"]:checked').each(function () { //遍历每一个名字为DelOpportCheck的复选框，其中选中的执行函数
            chk_value.push($(this).val()); //将选中的值添加到数组chk_value中      
        });
        OpportCtl.deleteOpport(chk_value)
    });

    Dom.on('click', '.alertOpportBtn', function () {
        $ModalChangeOpport.modal('show');
        var parent_dom = $(this).parents('tr');
        var Opport = new Object();
        Opport.title = parent_dom.find('.title').text();
        Opport.content = parent_dom.find('.content').text();
        Opport.entName = parent_dom.find('.entName').text();

        Opport.quota = parent_dom.find('.sales').text();
        Opport.opport_status = parent_dom.find('.opport_status').text();
        Opport.real_time = parent_dom.find('.real_time').text();
        OpportCtl.Opport.opportunities_id = parent_dom.find('.OpportId').data('id');
        
        OpportCtl.tags_html = parent_dom.find('.tag').html();
        $('#_tagAll .tag').removeClass('on');
        parent_dom.find('.tag').children('span').each(function () {
            var _class = ".tag_" + $(this).data('id');
            $('#_tagAll').find(_class).parent().addClass('on');
        });
		
		Dom.find('#_updateTags').html(OpportCtl.tags_html);

        Dom.find('#_OpportTitle').val(Opport.title);
        Dom.find('#_Content').val(Opport.content);
        Dom.find('#_Stage').find('option').each(function () {
            $(this).text() == Opport.opport_status && ($(this)[0].selected = true);
        });
        Dom.find('#_Type').find('option').each(function () {
            $(this).text() == Opport.type && ($(this)[0].selected = true);
        });
        Dom.find('#_OpportType').find('option').each(function () {
            $(this).text() == Opport.opportType && ($(this)[0].selected = true);
        })
        Dom.find('#_Quota').val(Opport.quota);
        Dom.find('#_RealTime').val(Opport.real_time);
    });

    $changeOpportBtn.on('click', function () {
        if (GlobalFunc.CheckForm($updateOpportForm)) {
            return
        };
        OpportCtl.changeOpport();
        OpportCtl.updateOpport(OpportCtl.Opport);
    });
    
    $closeOpportBtn.click(function () {
        !GlobalFunc.CheckChoose('DelOpportCheck') ? $ModalDelOpport.modal('show') : GlobalFunc.Alert.info('请勾选你需要关闭的商机！');
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
        OpportCtl.sortByKey(sort);
    });
    

    /*排序*/
//    Dom.on('click', '#timexSortBtn, #sortEstimateTimexBtn, #sortRealTimexBtn', function() {
//        var dom = $(this).parents('.panel-body').first();
//        var key = $(this).data('key');
//        var sortIndex = dom.find('.orderBy').find('input[type=radio]:checked').data('index');
//
//        OpportCtl.filter_info.sortIndex = sortIndex;
//        OpportCtl.filter_info.sortKey = key;
//        console.log(key)
//
//        if(key == 'timex') {
//            OpportCtl.filter_info.timex_begin = dom.find('.startValue').val();
//            OpportCtl.filter_info.timex_end = dom.find('.endValue').val();
//        } else if(key == 'estimate_time') {  
//            OpportCtl.filter_info.estimate_time_begin = dom.find('.startValue').val();
//            OpportCtl.filter_info.estimate_time_end = dom.find('.endValue').val();
//        } else if(key == 'real_time') {
//            OpportCtl.filter_info.real_time_begin = dom.find('.startValue').val();
//            OpportCtl.filter_info.real_time_end = dom.find('.endValue').val();
//        }
//
//        OpportCtl.bizFilter(OpportCtl.filter_info);
//    })
//
//    Dom.on('click', '#tagSortBtn', function () {
//        var dom = $(this).parents('.panel-body').first();
//        var tags = [];
//        dom.find('#allTagsBox').children('.tag.on').each(function () {
//            tags.push($(this).data('id'));
//        });
//        OpportCtl.filter_info.tags = tags;
//        OpportCtl.bizFilter(OpportCtl.filter_info);
//    });

    /*翻页*/
    $Opport_page.on('click', '.trun-left', function () {
        var page = OpportCtl.page_info.page - 1 || 1;
        OpportCtl.Flip(page);
    });
    $Opport_page.on('click', '.turn-right', function () {
        var page = OpportCtl.page_info.page = OpportCtl.page_info.page + 1 > OpportCtl.page_info.page_total ? OpportCtl.page_info.page_total : OpportCtl.page_info.page + 1;
        OpportCtl.Flip(page);
    });
    $Opport_page.on('click', '.branch-btn', function () {
        var page = Number($(this).text());
        OpportCtl.Flip(page);
    });


})($(document.body))

