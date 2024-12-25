var cdnDomain = "http://"+import.meta.env.VITE_BASE_URL + 'common' // 文件资源
/**
 * 动态表情
 */
const emojis = {
    '[桃心]': `<img width='30' src=${cdnDomain}/emojis/1.gif>`,
    '[拜托]': `<img width='30' src=${cdnDomain}/emojis/2.gif>`,
    '[差劲]': `<img width='30' src=${cdnDomain}/emojis/3.gif>`,
    '[闭嘴]': `<img width='30' src=${cdnDomain}/emojis/4.gif>`,
    '[真棒]': `<img width='30' src=${cdnDomain}/emojis/5.gif>`,
    '[头大]': `<img width='30' src=${cdnDomain}/emojis/6.gif>`,
    '[大哭]': `<img width='30' src=${cdnDomain}/emojis/7.gif>`,
    '[蛋糕]': `<img width='30' src=${cdnDomain}/emojis/8.gif>`,
    '[憨笑]': `<img width='30' src=${cdnDomain}/emojis/9.gif>`,
    '[恶魔]': `<img width='30' src=${cdnDomain}/emojis/10.gif>`,
    '[困]': `<img width='30' src=${cdnDomain}/emojis/11.gif>`,
    '[飞吻]': `<img width='30' src=${cdnDomain}/emojis/12.gif>`,
    '[奋斗]': `<img width='30' src=${cdnDomain}/emojis/13.gif>`,
    '[抓狂]': `<img width='30' src=${cdnDomain}/emojis/14.gif>`,
    '[干杯]': `<img width='30' src=${cdnDomain}/emojis/15.gif>`,
    '[勾引]': `<img width='30' src=${cdnDomain}/emojis/16.gif>`,
    '[购物]': `<img width='30' src=${cdnDomain}/emojis/17.gif>`,
    '[调皮]': `<img width='30' src=${cdnDomain}/emojis/18.gif>`,
    '[害羞]': `<img width='30' src=${cdnDomain}/emojis/19.gif>`,
    '[呲牙]': `<img width='30' src=${cdnDomain}/emojis/20.gif>`,
    '[鼓掌]': `<img width='30' src=${cdnDomain}/emojis/21.gif>`,
    '[无语]': `<img width='30' src=${cdnDomain}/emojis/22.gif>`,
    '[红包]': `<img width='30' src=${cdnDomain}/emojis/23.gif>`,
    '[骷髅]': `<img width='30' src=${cdnDomain}/emojis/24.gif>`,
    '[挥手]': `<img width='30' src=${cdnDomain}/emojis/25.gif>`,
    '[奖牌]': `<img width='30' src=${cdnDomain}/emojis/26.gif>`,
    '[拍手]': `<img width='30' src=${cdnDomain}/emojis/27.gif>`,
    '[惊恐]': `<img width='30' src=${cdnDomain}/emojis/28.gif>`,
    '[花痴]': `<img width='30' src=${cdnDomain}/emojis/29.gif>`,
    '[蜡烛]': `<img width='30' src=${cdnDomain}/emojis/30.gif>`,
    '[礼物]': `<img width='30' src=${cdnDomain}/emojis/31.gif>`,
    '[流汗]': `<img width='30' src=${cdnDomain}/emojis/32.gif>`,
    '[开心]': `<img width='30' src=${cdnDomain}/emojis/33.gif>`,
    '[亲亲]': `<img width='30' src=${cdnDomain}/emojis/34.gif>`,
    '[哭了]': `<img width='30' src=${cdnDomain}/emojis/35.gif>`,
    '[ok]': `<img width='30' src=${cdnDomain}/emojis/36.gif>`,
    '[吐了]': `<img width='30' src=${cdnDomain}/emojis/37.gif>`,
    '[庆祝]': `<img width='30' src=${cdnDomain}/emojis/38.gif>`,
    '[心碎]': `<img width='30' src=${cdnDomain}/emojis/39.gif>`,
    '[生病]': `<img width='30' src=${cdnDomain}/emojis/40.gif>`,
    '[气愤]': `<img width='30' src=${cdnDomain}/emojis/41.gif>`,
    '[可拍]': `<img width='30' src=${cdnDomain}/emojis/42.gif>`,
    '[难过]': `<img width='30' src=${cdnDomain}/emojis/43.gif>`,
    '[酷]': `<img width='30' src=${cdnDomain}/emojis/44.gif>`,
    '[睡觉]': `<img width='30' src=${cdnDomain}/emojis/45.gif>`,
    '[太阳]': `<img width='30' src=${cdnDomain}/emojis/46.gif>`,
    '[可爱]': `<img width='30' src=${cdnDomain}/emojis/47.gif>`,
    '[加油]': `<img width='30' src=${cdnDomain}/emojis/48.gif>`,
    '[握手]': `<img width='30' src=${cdnDomain}/emojis/49.gif>`,
    '[偷笑]': `<img width='30' src=${cdnDomain}/emojis/50.gif>`,
    '[玫瑰]': `<img width='30' src=${cdnDomain}/emojis/51.gif>`,
    '[大笑]': `<img width='30' src=${cdnDomain}/emojis/52.gif>`,
    '[晕]': `<img width='30' src=${cdnDomain}/emojis/53.gif>`,
    '[耶]': `<img width='30' src=${cdnDomain}/emojis/54.gif>`,
    '[惊讶]': `<img width='30' src=${cdnDomain}/emojis/55.gif>`,
    '[再见]': `<img width='30' src=${cdnDomain}/emojis/56.gif>`,
    '[招财猫]': `<img width='30' src=${cdnDomain}/emojis/57.gif>`,
    '[no]': `<img width='30' src=${cdnDomain}/emojis/58.gif>`,
    '[指]': `<img width='30' src=${cdnDomain}/emojis/59.gif>`,
    '[绅士]': `<img width='30' src=${cdnDomain}/emojis/60.gif>`
}

/**
 * 符号表情
 */
const symbol = [
    '😠', '😩', '😲', '😞', '😵', '😰', '😒', '😍', '😤', '😜', '😝', '😋', '😘', '😚', '😷',
    '😳', '😃', '😅', '😆', '😁', '😂', '😊', '☺', '😄', '😢',
    '😭', '😨', '😣', '😡', '😌', '😖', '😔', '😱', '😪', '😏', '😓', '😥', '😫', '😉',
    '✊', '✋', '✌', '👊', '👍', '☝', '👆', '👇', '👈', '👉',
    '👋', '👏', '👌', '👎'
]

const emojisKeys = Object.keys(emojis)

export const emojiList = {
    symbol,
    emojis
}

const regEmoji = emojisKeys.map((value) => '|\\' + value).join('').replace('|', '')

// 将字符串代码替换成 [XX] 字符
export function stringToCode(message) {
    if (message == null) {
        return ''
    }
    return message
        .replace('/::)', "[微笑]")
        .replace('/::~', "[撇嘴]")
        .replace('/::B', "[色]")
        .replace('/::|', "[发呆]")
        .replace('/::<', "[流泪]")
        .replace('/::$', "[害羞]")
        .replace('/::X', "[闭嘴]")
        .replace('/::Z', "[睡]")
        .replace('/::"(', "[大哭]")
        .replace('/::-|', "[尴尬]")
        .replace('/::@', "[发怒]")
        .replace('/::P', "[调皮]")
        .replace('/::D', "[呲牙]")
        .replace('/::O', "[惊讶]")
        .replace('/::(', "[难过]")
        .replace('/::-|', "[囧]")
        .replace('/::Q', "[抓狂]")
        .replace('/::T', "[吐]")
        .replace('/:,@P', "[偷笑]")
        .replace('/:,@-D', "[愉快]")
        .replace('/::d', "[白眼]")
        .replace('/:,@o', "[傲慢]")
        .replace('/::g', "[饥饿]")
        .replace('/:|-)', "[困]")
        .replace('/::!', "[惊恐]")
        .replace('/::L', "[流汗]")
        .replace('/::>', "[憨笑]")
        .replace('/::,@', "[悠闲]")
        .replace('/:,@f', "[奋斗]")
        .replace('/::-S', "[咒骂]")
        .replace('/:?', "[疑问]")
        .replace('/:,@x', "[嘘]")
        .replace('/:,@@', "[晕]")
        .replace('/::8', "[疯了]")
        .replace('/:,@!', "[衰]")
        .replace('/:xx', "[敲打]")
        .replace('/:bye', "[再见]")
        .replace('/:wipe', "[擦汗]")
        .replace('/:dig', "[抠鼻]")
        .replace('/:&-(', "[糗大了]")
        .replace('/:B-)', "[坏笑]")
        .replace('/:<@', "[左哼哼]")
        .replace('/:@>', "[右哼哼]")
        .replace('/::-O', "[哈欠]")
        .replace('/:>-|', "[鄙视]")
        .replace('/:P-(', "[委屈]")
        .replace('/::"|', "[快哭了]")
        .replace('/:X-)', "[阴险]")
        .replace('/::*', "[亲亲]")
        .replace('/:@x', "[吓]")
        .replace('/:8*', "[可怜]")
        .replace('/:8-', "[得意]")
        .replace('/::+', "[酷]")
        .replace('/:handclap', "[鼓掌]")
        .replace('/:hug', "[拥抱]")
        .replace('/:moon', "[月亮]")
        .replace('/:sun', "[太阳]")
        .replace('/:bome', "[炸弹]")
        .replace('/:!!!', "[骷髅]")
        .replace('/:pd', "[菜刀]")
        .replace('/:pig', "[猪头]")
        .replace('/:<W>', "[西瓜]")
        .replace('/:coffee', "[咖啡]")
        .replace('/:eat', "[饭]")
        .replace('/:heart', "[爱心]")
        .replace('/:break', "[心碎]")
        .replace('/:strong', "[强]")
        .replace('/:weak', "[弱]")
        .replace('/:share', "[握手]")
        .replace('/:v', "[胜利]")
        .replace('/:@)', "[抱拳]")
        .replace('/:jj', "[勾引]")
        .replace('/:ok', "[OK]")
        .replace('/:no', "[NO]")
        .replace('/:@@', "[拳头]")
        .replace('/:weak', "[差劲]")
        .replace('/:lvu', "[爱你]")
        .replace('/:showlove', "[嘴唇]")
        .replace('/:rose', "[玫瑰]")
        .replace('/:fade', "[凋谢]")
        .replace('/:beer', "[啤酒]")
        .replace('/:basketb', "[篮球]")
        .replace('/:oo', "[乒乓]")
        .replace('/:cake', "[蛋糕]")
        .replace('/:li', "[闪电]")
        .replace('/:footb', "[足球]")
        .replace('/:pd', "[刀]")
        .replace('/:kn', "[刀]")
        .replace('/:shit', "[便便]")
        .replace('/:ladybug', "[瓢虫]")
        .replace('/:gift', "[礼物]")
        .replace('/:jump', "[跳跳]")
        .replace('/:shake', "[发抖]")
        .replace('/:<O>', "[怄火]")
        .replace('/:circle', "[转圈]")
        .replace('/:kotow', "[磕头]")
        .replace('/:turn', "[回头]")
        .replace('/:skip', "[跳绳]")
        .replace('/:<&', "[左太极]")
        .replace('/:&>', "[右太极]")
        .replace('/:#-0', "[激动]")
        .replace('/:hiphot', "[乱舞]")
        .replace('/:kiss', "[献吻]")
        .replace('/:love', "[爱情]")
        .replace('/:<L>', "[飞吻]")
}


/**
 * 替换表情文字
 *
 * @param {String} content 需要替换的字符串
 */
export function textReplaceEmoji(content) {
    if (content.toLowerCase().indexOf("http://") >= 0) {
        content = `<a target="_blank" href="${content}">${content}</a>`;
    }
    if (content.toLowerCase().indexOf("https://") >= 0) {
        content = `<a target="_blank" href="${content}">${content}</a>`;
    }
    //替换 www. bbs. 等开头网址
    if (content.toLowerCase().indexOf("www.") >= 0 || content.toLowerCase().indexOf("bbs.") >= 0) {
        content = content.replace(/(^|[^\/\\\w\=])((www|bbs)\.(\w)+\.([\w\/\\\.\=\?\+\-~`@\'!%#]|(&amp;))+)/g, "$1<a target=\"_blank\" href=http://$2>$2</a>");
    }
    content = stringToCode(content)
    return content.replace(new RegExp(`(${regEmoji})`, 'gi'), ($0, $1) => {
        return emojis[$1]
    })
}