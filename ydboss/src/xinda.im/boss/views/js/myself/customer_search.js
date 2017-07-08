(function ($) {
    "user struct";
    /**
     * Statement
     */
    var GlobalFunc = new Global(),
        CustomerCtl = new Customer(),
        Tpl_Customer_list = _.template($('#TplCustomerList').html().replace(/\&lt\;/ig, '\<')),
        $Customer_page = $('#Customer_page'),
        $queryBtn = $('#queryBtn'),
        $keyWord = $('#keyWord'),
        $provincetagList = $('#provincetagList'),
        $TagtagList = $('#TagtagList'),
        $staffFrom = $('#staffFrom'),
        $staffTo = $('#staffTo'),
        $timexFrom = $('#timexFrom'),
        $timexTo = $('#timexTo'),
        $followTimeFrom = $('#followTimeFrom'),
        $CustomerListBox = $('#CustomerListBox'),
        $followTimeTo = $('#followTimeTo'),
        $loadMask = $('#loadMask')

    var That;

    /**
     * Customer [CLASS]
     */
    function Customer() {
        That = this;

        this.getType = function (href) {
            return {
                "all_customer_search": "all",
                "agent_customer_search": "agent"
            }[href] || "sale";
        };

        this.type = this.getType(location.href.match(/.*\/(.+).html/)[1]);

        this.queryBefore = function () {
            var o = {
                width: $CustomerListBox.parent().width(),
                height: $CustomerListBox.parent().height(),
                lineHeight: $CustomerListBox.parent().height()
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

        this.showCustomerList = function (data) {
            data.link = {
                "all": "get_super_customer_detail",
                "agent": "get_agent_customer_detail"
            }[this.type] || "get_emp_customer_detail";
            
            var h_table = Tpl_Customer_list(data);
            $CustomerListBox.html(h_table);
            $Customer_page.find('.text-primary').text(data.page_info.record);
            $('.popovers').popover({
                html: true
            });
            this.queryAfter();
        };
    };

    /**
     * Customer Prototype
     */
    Customer.prototype.getCustomerList = function () {
        var query = {
            "all": "/customer.list",
            "agent": "/agent.customer.list"
        }[this.type] || "/employee.customer.list";

        $.post(query, {
            page: 1
        }, function (data) {
            if (!!data && data.code == 0) {
                That.showCustomerList(data);
            }
        })
    };

    Customer.prototype.screenCustomerList = function (o) {
        if (JSON.stringify(o) == "{}") {
            this.getCustomerList();
            return;
        }
        var query = {
            "all": "/customer_super_filter",
            "agent": "/customer_agent_filter"
        }[this.type] || "/customer_emp_filter";

        $.post(query, o, function (data) {
            if (!!data && data.code == 0) {
                That.showCustomerList(data);
            }
        })
    }

    CustomerCtl.getCustomerList();

    /**
     * All Event Bind 
     */
    $queryBtn.click(function () {
        var o = new Object(),
            $this = $(this),
            provinces = [],
            tags = [];
        $provincetagList.children('.tag_li').each(function () {
            provinces.push($(this).data('id'));
        });
        $TagtagList.children('.tag_li').each(function () {
            tags.push($(this).data('id'));
        });

        CustomerCtl.queryBefore();


        !!$.trim($keyWord.val()) && (o.keyword = $keyWord.val());
        provinces.length > 0 && (o.provinces = provinces);
        tags.length > 0 && (o.tags = tags);
        !!$.trim($staffFrom.val()) && (o.emp_count_from = $staffFrom.val());
        !!$.trim($staffTo.val()) && (o.emp_count_to = $staffTo.val());
        !!$.trim($timexFrom.val()) && (o.timex_from = $timexFrom.val());
        !!$.trim($timexTo.val()) && (o.timex_to = $timexTo.val());
        !!$.trim($followTimeFrom.val()) && (o.follow_time_from = $followTimeFrom.val());
        !!$.trim($followTimeTo.val()) && (o.follow_time_to = $followTimeTo.val());

        CustomerCtl.screenCustomerList(o);
    });


})(jQuery)