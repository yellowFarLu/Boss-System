(function ($) {
    "user struct";
    /**
     * Statement
     */
    var GlobalFunc = new Global(),
        OpportCtl = new Opport(),
        Tpl_opport_list = _.template($('#TplOpportList').html().replace(/\&lt\;/ig, '\<')),
        $opport_page = $('#opport_page'),
        $queryBtn = $('#queryBtn'),
        $keyWord = $('#keyWord'),
        $TagtagList = $('#TagtagList'),
        $salesFrom = $('#salesFrom'),
        $salesTo = $('#salesTo'),
        $timexFrom = $('#timexFrom'),
        $timexTo = $('#timexTo'),
        $estimateTimeFrom = $('#estimateTimeFrom'),
        $estimateTimeTo = $('#estimateTimeTo'),
        $realTimeFrom = $('#realTimeFrom'),
        $realTimeTo = $('#realTimeTo'),
        $opportListBox = $('#opportListBox'),
        $loadMask = $('#loadMask')

    var That;

    /**
     * Opport [CLASS]
     */
    function Opport() {
        That = this;

        this.getType = function (href) {
            return {
                "all_opport_search": "all",
                "agent_opport_search": "agent"
            }[href] || "sale";
        };

        this.type = this.getType(location.href.match(/.*\/(.+).html/)[1]);

        this.queryBefore = function () {
            var o = {
                width: $opportListBox.parent().width(),
                height: $opportListBox.parent().height(),
                lineHeight: $opportListBox.parent().height()
            }
            $loadMask.css({
                "width": o.width,
                "height": o.height,
                "lineHeight": o.lineHeight + "px"
            }).show();
            $queryBtn.css("opacity", ".5");
        };
        this.queryAfter = function () {
            $loadMask.hide();
            $queryBtn.css("opacity", "1");
        };

        this.showOpportList = function (data) {
            data.link = {
                "all": "get_super_biz_detail",
                "agent": "get_agent_biz_detail"
            }[this.type] || "get_emp_biz_detail";
            
            _.map(data.list, function (v) {
                v.status = GlobalFunc.OpportStage[v.status];
            });
            
            var h_table = Tpl_opport_list(data);
            $opportListBox.html(h_table);
            $opport_page.find('.text-primary').text(data.page_info.record);
            this.queryAfter();
        };
    };

    /**
     * Opport Prototype
     */
    Opport.prototype.getOpportList = function () {
        var query = {
            "all": "/all.opport.list",
            "agent": "/agent_opport_list"
        }[this.type] || "/opport.list";

        $.post(query, {
            page: 1
        }, function (data) {
            if (!!data && data.code == 0) {
                That.showOpportList(data);
            }
        })
    };

    Opport.prototype.screenOpportList = function (o) {
        if (JSON.stringify(o) == "{}") {
            this.getOpportList();
            return;
        }
        var query = {
            "all": "/opport_super_filter",
            "agent": "/opport_agent_filter"
        }[this.type] || "/opport_emp_filter";

        $.post(query, o, function (data) {
            if (!!data && data.code == 0) {
                That.showOpportList(data);
            }
        })
    }

    OpportCtl.getOpportList();

    /**
     * All Event Bind 
     */
    $queryBtn.click(function () {
        var o = new Object(),
            $this = $(this),
            tags = [];
        $TagtagList.children('.tag_li').each(function () {
            tags.push($(this).data('id'));
        });

        OpportCtl.queryBefore();


        !!$.trim($keyWord.val()) && (o.keyword = $keyWord.val());
        tags.length > 0 && (o.tags = tags);
        !!$.trim($salesFrom.val()) && (o.sales_from = $salesFrom.val());
        !!$.trim($salesTo.val()) && (o.sales_to = $salesTo.val());
        !!$.trim($timexFrom.val()) && (o.timex_from = $timexFrom.val());
        !!$.trim($timexTo.val()) && (o.timex_to = $timexTo.val());
        !!$.trim($estimateTimeFrom.val()) && (o.estimate_time_from = $estimateTimeFrom.val());
        !!$.trim($estimateTimeTo.val()) && (o.estimate_time_to = $estimateTimeTo.val());
        !!$.trim($realTimeFrom.val()) && (o.real_time_from = $realTimeFrom.val());
        !!$.trim($realTimeTo.val()) && (o.real_time_to = $realTimeTo.val());

        OpportCtl.screenOpportList(o);
    });


})(jQuery)