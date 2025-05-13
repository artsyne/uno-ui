(function ($) {
    // 定义 jQuery 插件 biurenGrid
    $.fn.biurenGrid = function (options) {
        // 初始化 HTML 数组
        const aryHtml = [];

        // 添加表格的起始标签
        aryHtml.push(`<table id="${options.id}" class="${options.className}">`);

        // 如果有表格标题，则添加标题
        if (options.caption) {
            aryHtml.push(`<caption>${options.caption}</caption>`);
        }

        // 如果有表头模板，则添加表头
        if (options.template.head) {
            aryHtml.push("<thead><tr>");
            options.template.head.forEach(head => {
                // 获取表头文本
                const title = head.text || head;
                aryHtml.push(`<td>${title}</td>`);
            });
            aryHtml.push("</tr></thead>");
        }

        // 添加表格主体的起始标签
        aryHtml.push("<tbody>");
        const tempRow = options.template.row;
        const ds = options.ds;

        // 如果数据源为空，则添加无数据提示
        if (ds.length === 0) {
            aryHtml.push(`<tr><td>${options.noData}</td></tr>`);
        } else {
            // 遍历数据源，生成表格行
            ds.forEach((data, i) => {
                // 每列数目达到指定值时，换行
                if (i % options.col === 0) {
                    if (i !== 0) aryHtml.push("</tr>");
                    aryHtml.push("<tr");

                    // 如果有行样式，则添加样式
                    if (tempRow.className) {
                        let temp = tempRow.className;
                        // 如果有交替行样式，则应用交替样式
                        if (tempRow.switchClassName && i % 2 !== 0) {
                            temp = tempRow.switchClassName;
                        }
                        aryHtml.push(` class="${temp}"`);
                    }

                    // 如果有行标记，则添加标记
                    if (tempRow.mark) {
                        aryHtml.push(` ${tempRow.mark}`);
                    }

                    // 如果有行事件，则添加事件
                    if (tempRow.event) {
                        tempRow.event.forEach(event => {
                            aryHtml.push(` ${event.name}="${event.func}"`);
                        });
                    }

                    aryHtml.push(">");
                }

                const tempCell = options.template.cell;
                let cellMark = "";

                // 如果有单元格属性，则添加单元格事件
                if (options.template.cellAttr && options.template.cellAttr.event) {
                    options.template.cellAttr.event.forEach(event => {
                        cellMark += ` ${event.name}="${event.func}"`;
                    });
                }

                // 遍历单元格模板，生成单元格
                tempCell.forEach(cell => {
                    aryHtml.push(`<td ${cellMark}`);
                    if (cell.className) {
                        aryHtml.push(` class="${cell.className}"`);
                    }
                    aryHtml.push(`>${replaceParam(cell.text, getValueFromAry(data, cell.value))}</td>`);
                });
            });
        }

        // 结束表格主体标签
        aryHtml.push("</tr></tbody>");

        // 如果有表尾模板且数据源不为空，则添加表尾
        if (options.template.foot && ds.length !== 0) {
            aryHtml.push(`<tfoot><tr><td colspan="${options.template.cell.length * options.col}">`);
            let temp = options.template.foot.text;
            if (options.pageLink) {
                temp = temp.replace("{0}", `<div class="biurennav">${options.pageLink}</div>`);
            }
            aryHtml.push(`<div class="pages">${temp}</div>`);
            aryHtml.push(`<input type="hidden" name="${options.id}ViewState" id="${options.id}ViewState}" />`);
            aryHtml.push("</td></tr></tfoot>");
        }

        // 结束表格标签
        aryHtml.push("</table>");
        // 将生成的 HTML 设置到容器中
        this.html(aryHtml.join(""));

        return this;
    };

    // 替换参数占位符的函数
    function replaceParam(text, value) {
        return value.reduce((result, val, i) => {
            const reg = new RegExp(`{(${i})}`, "g");
            return result.replace(reg, val);
        }, text);
    }

    // 从数组中获取值的函数
    function getValueFromAry(ary, value) {
        return value.map(val => {
            let temp = ary[val.pos || val] || val.def || "";
            temp = temp.toString();

            if (val.parse && val.parse[temp]) {
                temp = val.parse[temp];
            }

            if (val.func) {
                temp = eval(val.func)(temp);
            }

            if (val.maxlen) {
                temp = truncate(temp, val.maxlen, "...");
            }

            if (!val.isHtml) {
                temp = $.HtmlEncode(temp);
            }

            return temp;
        });
    }

    // 截断字符串的函数
    function truncate(val, num, suffix) {
        if (val.replace(/[^\x00-\xff]/g, "**").length <= num) {
            return val;
        }
        let bytenum = 0;
        let rs = "";
        const pt = /[^\x00-\xff]/;
        num -= suffix.length;

        for (let i = 0; i < num; i++) {
            const ch = val.substr(i, 1);
            bytenum += ch.match(pt) ? 2 : 1;
            if (bytenum > num) break;
            rs += ch;
        }

        return rs + suffix;
    }
})(jQuery);