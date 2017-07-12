$(document).ready(function(){
$("h2,h3").each(function(i,item){
var tag = $(item).get(0).localName;
$(item).attr("id","wow"+i);
$("#category").append('<a class="new'+tag+'" href="#wow'+i+'">'+$(this).text()+'</a></br>');
$(".newh2").css("margin-left",20);
$(".newh3").css("margin-left",40);
});
});