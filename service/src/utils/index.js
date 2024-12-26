import axios from 'axios'
import { compress, compressAccurately } from 'image-conversion';
import { message } from 'ant-design-vue';

/**
 * 全局方法
 * @localGet  ---获取缓存
 * @localSet  ---存入缓存
 * @localRemove  ---清空所有缓存
 * @localEmpty  ---删除缓存
 * @beforeUpload  ---图片文件前
 * @videoBeforeUpload ---上传视频前
 * @uploadFile  ---上传文件
 * @FileBase64  ---文件转base64
 * @accDiv  ---除法
 * @accMul  ---乘法
 * @accAdd  ---加法
 * @accSub  ---减法
 * 
 */

/**
 * 获取缓存
 * @param {缓存名} key 
 * @returns 
 */
export function localGet(key) {
    const value = window.localStorage.getItem(key)
    try {
        return JSON.parse(window.localStorage.getItem(key))
    } catch (error) {
        return value
    }
}
/**
 * 存入缓存
 * @param {缓存名} key 
 * @param {缓存值} value 
 */
export function localSet(key, value) {
    window.localStorage.setItem(key, JSON.stringify(value))
}
/**
 * 删除缓存
 * @param {缓存名} key 
 */
export function localEmpty(key) {
    window.localStorage.removeItem(key)
}
/**
 * 清空所有缓存
 */
export function localRemove() {
    window.localStorage.clear()

}
/**
 * 上传图片前
 * @param {文件} file 
 * @returns 
 */
export async function beforeUpload(file) {
    return new Promise((resolve, reject) => {
        const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
        if (!isJpgOrPng) {
            return message.error('只能上传图片!');
        }else{
            if (file.size / 1024 > 200) compressAccurately(file, 200).then(res => { resolve(res) }) // 大于 200 kb 就压缩
            else resolve(file)
        }
    })
}
/**
 * 上传视频前
 * @param {文件} file 
 * @returns 
 */
export async function videoBeforeUpload(file) {
    const isJpgOrPng = file.type === 'video/mp4';
    if (!isJpgOrPng) {
        message.error('只能上传视频!');
    }
    return isJpgOrPng
}
/**
 * 上传文件
 * @param {文件} file 
 * @returns 
 */
// export async function uploadFile(file) {
//     return new Promise(async (resolve, reject) => {
//         //获取上传签名
//         const res = await axios.post("/common/api/oss_config", { file_name: file.name }, false);
//         var xhr = new XMLHttpRequest();
//         xhr.open('PUT', res.data.url, true);
//         xhr.onload = function (e) {
//             resolve(res.data.file_name)
//         };
//         xhr.send(file); //file 是要上传的文件对象
//     })
// }
export async function uploadFile(file) {
    return new Promise((resolve, reject) => {
        let formData = new FormData();
        formData.append('image', file);
        console.log(import.meta.env.VITE_BASE_URL)
        axios({
            method: 'post',
            url: "http://"+import.meta.env.VITE_BASE_URL + 'api/system/upload',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            data: formData,
        }).then(({ data: res }) => {
            resolve(res.data.file_name)
        }).catch(err => { });
    })
}
/**
 * 文件转base64
 * @param {文件} file 
 * @returns 
 */
export async function FileBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onloadend = (e) => {
            resolve(e.target.result)
        }
    })    
}
/**
 * 标准时间转换
 * @param {标准时间} time 
 * @returns 
 */
export async function getTime(time) {
    var date = new Date(time)
    var y = date.getFullYear()
    var m = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1)
    var d = (date.getDate() < 10 ? '0' + (date.getDate()) : date.getDate())
    var hour = (date.getHours() < 10 ? '0' + (date.getHours()) : date.getHours())
    var minute = (date.getMinutes() < 10 ? '0' + (date.getMinutes()) : date.getMinutes())
    return y + '-' + m + '-' + d + ' ' + hour + ":" + minute
}
/**
 * 除法
 * @param {除数}} arg1 
 * @param {被除数} arg2 
 * @returns 
 */
function accDiv(arg1, arg2) {
    var t1 = 0,
        t2 = 0,
        r1, r2;
    try {
        t1 = arg1.toString().split(".")[1].length
    } catch (e) {}
    try {
        t2 = arg2.toString().split(".")[1].length
    } catch (e) {}
    r1 = Number(arg1.toString().replace(".", ""));
    r2 = Number(arg2.toString().replace(".", ""));
    if (r2 == 0) {
        return Infinity;
    } else {
        return (r1 / r2) * Math.pow(10, t2 - t1);
    }
}
/**
 * 乘法
 * @param {乘数} arg1 
 * @param {乘数} arg2 
 * @returns 
 */
function accMul(arg1, arg2) {
    var m = 0,
        s1 = arg1.toString(),
        s2 = arg2.toString();
    try {
        m += s1.split(".")[1].length
    } catch (e) {}
    try {
        m += s2.split(".")[1].length
    } catch (e) {}

    return Number(s1.replace(".", "")) * Number(s2.replace(".", "")) / Math.pow(10, m);
}
/**
 * 加法
 * @param {加数} arg1 
 * @param {加数} arg2 
 * @returns 
 */
function accAdd(arg1, arg2) {
    var r1, r2, m;
    try {
        r1 = arg1.toString().split(".")[1].length;
    } catch (e) {
        r1 = 0;
    }
    try {
        r2 = arg2.toString().split(".")[1].length;
    } catch (e) {
        r2 = 0;
    }
    m = Math.pow(10, Math.max(r1, r2));
    return (arg1 * m + arg2 * m) / m;
}
/**
 * 减法
 * @param {减数} arg1 
 * @param {被减数} arg2 
 * @returns 
 */
function accSub(arg1, arg2) {
    var r1, r2, m, n;
    try {
        r1 = arg1.toString().split(".")[1].length;
    } catch (e) {
        r1 = 0;
    }
    try {
        r2 = arg2.toString().split(".")[1].length;
    } catch (e) {
        r2 = 0;
    }
    m = Math.pow(10, Math.max(r1, r2));
    //动态控制精度长度
    n = (r1 >= r2) ? r1 : r2;
    return ((arg2 * m - arg1 * m) / m).toFixed(n);
}
//加法。
Number.prototype.add = function(arg) {
    return accAdd(arg, this);
};
//减法
Number.prototype.sub = function(arg) {
    return accSub(arg, this);
};
//乘法
Number.prototype.mul = function(arg) {
    return accMul(arg, this);
};
//除法
Number.prototype.div = function(arg) {
    return accDiv(this, arg);
};
export {
    accDiv,
    accMul,
    accAdd,
    accSub,
}