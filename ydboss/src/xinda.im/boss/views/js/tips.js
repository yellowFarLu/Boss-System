
/**
 * 1、提示插件
 * 2、存正则表达式变量
 */ 

/* 
*提示插件 
*使用方法:showTip('提示内容',3) 
*/  
  
var tip = null;  
  
//显示提示文本，text为要显示的文本，second表示自动消失的秒数  
function showTip(text,second){  
      
    if(!tip){  
        tip = document.createElement('div');  
        tip.style.backgroundColor = '#000';  
        tip.style.color = '#fff';  
        tip.style.fontSize = '18px';  
        tip.style.fontWeight = 'bold';  
        tip.style.fontFamily = '微软雅黑';  
        tip.style.width = 'auto';  
        tip.style.height = 'auto';  
        tip.style.padding = '10px 30px 10px 30px';  
        tip.style.borderRadius = '6px';  
        tip.style.opacity = '0.7';  
        tip.style.position = 'fixed';  
        tip.style.zIndex = '99999';  
        document.body.appendChild(tip);  
    }  
    tip.innerHTML = text;  
    tip.style.display = 'inline-block';  
  
    tip.style.top = (document.documentElement.clientHeight/2 - tip.offsetHeight/2) + 'px';  
    tip.style.left = (document.documentElement.clientWidth/2 - tip.offsetWidth/2) + 'px';  
      
    setTimeout(function disappear(){  
        tip.style.display = 'none';  
    },second*1000);  
      
}


/**
 * 正则表达式变量
 */
var PhoneRegex = eval("/^1\d{10}$/");
var EmailRegex = eval("/^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[a-zA-Z]{2,3}$/");

