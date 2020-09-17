import Vue from 'vue'
import VueI18n from 'vue-i18n'

Vue.use(VueI18n)

const DEFAULT_LANG = 'zh'
const LOCALE_KEY = 'i18n'

const locales = {
    zh: require('./zh-CN.json'),
    en: require('./en-US.json'),
}

const i18n = new VueI18n({
    locale: DEFAULT_LANG,
    messages: locales,
})

export const setup = lang => {
    if (lang === undefined) {
        lang = window.localStorage.getItem(LOCALE_KEY)
        if (locales[lang] === undefined) {
            lang = DEFAULT_LANG
        }
    }
    window.localStorage.setItem(LOCALE_KEY, lang)

    Object.keys(locales).forEach(lang => {
        document.body.classList.remove(`lang-${lang}`)
    })
    document.body.classList.add(`lang-${lang}`)
    document.body.setAttribute('lang', lang)

    Vue.config.lang = lang
    i18n.locale = lang
}

setup()
export default i18n