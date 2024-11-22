import axios from "@/utils/axios"
// import OSS from "ali-oss"
import mapList from "@/utils/map";

/**
 * 全局方法
 * @localGet  ---获取缓存
 * @localSet  ---存入缓存
 * @localRemove  ---清空所有缓存
 * @localEmpty  ---删除缓存
 * @getElementData  ---获取元素数据（高度、宽度）
 * @beforeUpload  ---上传文件前
 * @uploadImg  ---oss上传图片
 * @getAddressOptions  ---三级联动获取地址
 * @addressContext  ---根据地址ID获取地址名称
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
 * 获取元素数据（高度、宽度）
 * @param {元素id} id 
 * @param {获取数据类型 h:高 w:宽} type 
 * @returns 
 */
export function getElementData(id, type) {
    if (type == 'h') return document.getElementById(id).clientHeight
    else document.getElementById(id).clientWidth
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
        message.error('图片大小必须小于2MB!');
    }
    return isJpgOrPng && isLt2M;
}
/**
 * oss上传图片
 * @param {文件} file 
 * @returns 
 */
export async function uploadImg(file) {
    const res = await axios.get('/oss/policy');
    let client = new OSS({
        region: res.region,
        accessKeyId: res.accessKeyId,
        accessKeySecret: res.accessKeySecret,
        bucket: res.bucket,
    });
    var fileName = `${Date.parse(new Date())}` + file.name;
    let host = res.host + fileName
    const result = await client.put(host, file);
    return result.url
}
/**
 * 三级联动获取地址
 * @param {层级数} level 
 * @param {父级id} id 
 * @returns 
 */
export async function getAddressOptions(level, id = 0) {
    let data = []
    if (level == 1) {
        mapList.map.forEach(element => {
            if (element.id == id) {
                data = element.children
            }
        });
    } else if (level == 2) {
        mapList.map.forEach(element => {
            element.children.forEach(city => {
                if (city.id == id) data = city.children
            });
        });
    } else {
        mapList.map.forEach(element => {
            data.push(element)
        });
    }
    return data
}
/**
 * 根据地址ID获取地址名称
 * @param {地址id数组} e 
 * @returns 
 */
export function addressContext(e) {
    let obj = ""
    let arry, arry1, arry2 = []
    if (e[0]) {
        arry = mapList.map.filter((i) => i.id == e[0]);
        obj = arry[0].name
    } else if (e[1]) {
        arry1 = arry[0].children.filter((i) => i.id == e[1])
        obj = obj + ' - ' + arry1[0].name
    } else if (e[2]) {
        arry2 = arry1[0].children.filter((i) => i.id == e[2])
        obj = obj + ' - ' + arry2[0].name
    }
    return obj
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