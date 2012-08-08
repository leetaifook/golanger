$(function(){
    function parse_bytes(a,b){a=undefined===a?0:a,b=undefined===b?1024:b;var c=b,d=b*c,e=b*d;return a=Number(a),a>c?a>d?a=a>e?(a/e).toFixed(2).toString()+" GB":(a/d).toFixed(2).toString()+" MB":a=(a/c).toFixed(2).toString()+" KB":a>0&&(a+=" B"),a}

    $(".nav-collapse > ul > li > a[href='"+document.location.pathname+"']").parent("li").addClass("active");
});