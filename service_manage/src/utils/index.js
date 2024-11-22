import axios from 'axios';
import { message } from 'ant-design-vue';

/**
 * 全局方法
 * @localGet  ---获取缓存
 * @localSet  ---存入缓存
 * @localRemove  ---清空所有缓存
 * @localEmpty  ---删除缓存
 * @beforeUpload  ---上传文件前
 * @uploadFile  ---上传文件
 * @getElementData  ---获取元素数据（高度、宽度）
 * @getTime ---时间转换
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
 * 上传文件前
 * @param {文件} file 
 * @returns 
 */
export async function beforeUpload(file) {
    const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
    if (!isJpgOrPng) {
        message.error('只能上传图片!');
    }
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
        message.error('文件大小必须小于2MB!');
    }
    return isJpgOrPng && isLt2M;
}
/**
 * 上传图片
 * @param {文件} file 
 * @returns 
 */
export function uploadFile(file) {
    return new Promise((resolve, reject) => {
        let formData = new FormData();
        formData.append('image', file);
        axios({
            method: 'post',
            url: import.meta.env.VITE_BASE_URL + '/api/system/upload',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            data: formData,
        }).then(({ data: res }) => {
            resolve(res.data.file_name)
        }).catch(err => { });
    })
}
/**
 * 获取元素数据（高度、宽度）
 * @param {元素id} id 
 * @param {获取数据类型 h:高 w:宽} type 
 * @returns 
 */
export function getElementData (id,type) {
    if (type == 'h') return document.getElementById(id).clientHeight
    else document.getElementById(id).clientWidth
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
    return y + '-' + m + '-' + d
}
/**
 * 获取当前时间
 * @returns 
 */
export async function getNowTime() {
    var date = new Date();
    let year = date.getFullYear();
    let month = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1);
    let day = (date.getDate() < 10 ? '0' + (date.getDate()) : date.getDate());
    let hour = date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
    let minute = date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
    let second = date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds(); 
    let weeks = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];
    return year + '-' + month + '-' + day + ' ' + weeks[date.getDay()] + ' ' +  + hour + ':' + minute + ':' + second;
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