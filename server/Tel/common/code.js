const key = "USER_LOGIN_TOKEN_CACHE";


// 基本链接
const base = "/user/auth/local_storage"



let uuid = localStorage.getItem(key)
console.log("设置前UUID",uuid)

// 如果没有uuid则生成一个uuid
if (!uuid || uuid === "") {
    uuid = CreateUUID()
    console.log("设置UUID",uuid)
    localStorage.setItem(key, uuid);
}

const url = document.getElementById("#imagesUrl").src
document.getElementById("#imagesUrl").src =  url + "&uid="+uuid

//
// console.log("code,",code)
// console.log("baseLink,",baseLink)
// console.log("bindLink,",bindLink)
// console.log("newUUid,",newUUid)
// console.log("actionUrl,",actionUrl)
//
// console.log("join",CreateLink(uuid))
//
//
// let currentUuid = localStorage.getItem(key)
//
// // 如果uuid不同则绑定uuid并重定向
// if (currentUuid && currentUuid !== "") {
//     // alert("transfer 不是空 uuid 準備綁定 currentUuid "+currentUuid + "   --  newUUid"+newUUid)
//     bindUuIDAndLocation(newUUid, currentUuid,"transfer")
// }else{
//     // 如果本地没有储存uuid则储存一个uuid
//     if (!currentUuid || currentUuid === "" ) {
//         localStorage.setItem(key, newUUid);
//     }
//
//     const url = document.getElementById("#imagesUrl").src
//     document.getElementById("#imagesUrl").src =  url + "&uid="+currentUuid
//     // 带着uuid跳转到目的地
//     window.location.href = CreateLink(newUUid)
// }
//
// // 生成一个UUID
// function CreateUUID() {
//     function guid2() {
//         function S4() {
//             return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
//         }
//
//         return code + "_" + (S4() + S4() + "_" + S4() + "_" + S4() + "_" + S4() + "_" + S4() + S4() + S4());
//     }
//
//     return guid2() + "_" + new Date().getTime()
// }
//
// // 根据uuid创建跳转连接
// function CreateLink(uuid) {
//     return baseLink + code + "/" + uuid
// }

// // 绑定uuid并重定向
// function bindUuIDAndLocation(newUuid,currentUuid,action) {
//     const xhr = new XMLHttpRequest();
//     xhr.open('POST', bindLink+"/bind_uuid/"+newUuid+"/"+currentUuid, true);
//     xhr.onload = function (e) {
//         console.log(e)
//         // window.location.href = CreateLink(currentUuid)
//         const url = document.getElementById("#imagesUrl").src
//         document.getElementById("#imagesUrl").src =  url + "&uid="+currentUuid
//     };
//     xhr.send(JSON.stringify({action:action})); //file 是要上传的文件对象
// }
