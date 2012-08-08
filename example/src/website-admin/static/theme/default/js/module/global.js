$(function () {
    $("#btn_cancel").bind("click", function(){
        window.location = "/module/";
    })
    $("#btn_cancel_create").bind("click", function(){
        window.location = "/module/";
    })
})

var ck_modulename = /^[A-Za-z0-9_]{1,20}$/;


function validateModulename(modulename) {
    if (!ck_modulename.test(modulename)) {
        return false;
    }
    return true;
}