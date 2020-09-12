<template>
  <div class="home">
    <check />
    <div v-if="!showPost" id="list">
      <label
        class="tab"
        :class="showRead && showUnread ? 'select' : ''"
        @click="
          showUnread = true;
          showRead = true;
        "
        >All</label
      >
      |
      <label
        class="tab"
        :class="showRead && !showUnread ? 'select' : ''"
        @click="
          showUnread = false;
          showRead = true;
        "
        >Read</label
      >
      |
      <label
        class="tab"
        :class="showUnread && !showRead ? 'select' : ''"
        @click="
          showUnread = true;
          showRead = false;
        "
        >Unread</label
      >
      <div v-for="(i, index) in post" :key="index" style="text-align: left">
        <div
          class="post"
          :class="i.read ? 'read' : 'unread'"
          v-if="
            (showRead && i.read && !showUnread) ||
              (showUnread && !i.read) ||
              (showRead && showUnread)
          "
        >
          <a style="font-size: small;color: rgba(0,0,0,.7)"
            >{{ i.source }} >>
          </a>
          <a
            style="font-size: large"
            @click="
              nowPost = index;
              showPost = true;
              i.read = false;
              change(index);
            "
            >{{ i.title }}
          </a>
          <b-icon-check-square-fill
            style="float: right;margin: 5px;color: rgb(69,123,48)"
            v-if="i.read"
            @click="change(index)"
            >read</b-icon-check-square-fill
          >
          <b-icon-check-square
            style="float: right;margin: 5px"
            v-else
            @click="change(index)"
            >unread</b-icon-check-square
          >
          <a
            style="font-size: small;color: rgba(0,0,0,.7);float: right;margin: 5px"
            class="postdate"
          >
            {{ i.date }}
          </a>
        </div>
      </div>
    </div>
    <div v-else id="postinfo">
      <label
        class="tab"
        style="margin: 5px;float: left"
        @click="showPost = false"
        >Back</label
      >
      <div style="float: right">
        <b-icon-check-square-fill
          style="float: right;margin: 5px;color: rgb(69,123,48)"
          v-if="post[nowPost].read"
          @click="change(nowPost)"
          >read</b-icon-check-square-fill
        ><label v-if="post[nowPost].read">已读</label>
        <b-icon-check-square
          style="float: right;margin: 5px"
          v-if="!post[nowPost].read"
          @click="change(nowPost)"
          >unread</b-icon-check-square
        ><label v-if="!post[nowPost].read">未读</label>
      </div>

      <h1>{{ post[nowPost].title }}</h1>
      <span>{{ post[nowPost].source }}</span> |
      <a :href="post[nowPost].link">Link</a>
      <b-card id="postcontext" style="margin: 15px" v-html="post[nowPost].desc">
      </b-card>
    </div>
  </div>
</template>

<script>
import check from "@/components/check";
export default {
  name: "Overview",
  components: {
    check
  },
  methods: {
    change: function(index) {
      this.post[index].read ? this.unread(index) : this.read(index);
      this.post[index].read = this.post[index].read ? false : true;
    },
    read: function(index) {
      console.log(index);
    },
    unread: function(index) {
      console.log(index);
    },
    getReadList: function() {
      return;
    },
    getPostList: function() {
      return;
    },
    getPostContent: function() {
      return;
    }
  },
  data() {
    return {
      post: [
        {
          title: "liunx jetbrains 软件输入中文",
          date: "2020-06-13 13:00",
          link: "https://blog.lvcshu.com",
          source: "blog.lvcshu.com",
          desc: "这里是描述",
          read: false
        },
        {
          title: "git 同步上游代码",
          date: "2020-05-31 13:00",
          link: "https://blog.lvcshu.com",
          source: "blog.lvcshu.com",
          desc:
            '<p>好久没有更新博客啦，上来写点碎碎念<a id="more"></a></p><h1 id="博客主题更新"><a href="#博客主题更新" class="headerlink" title="博客主题更新"></a>博客主题更新</h1><p>最近在上网课之余腾出手来将 Hexo 主题做了一点点更新，主要是将<a href="https://blog.lvcshu.com">入口页面</a>做的漂亮了一点，加上了背景图片以及把顶栏 CSS 调整成了透明来更加适应图片背景。</p><p>然后发现如果顶栏一直透明的话滚动到了文章列表会比较难看，就加了一点 js 使顶栏能自己切换透明以及白色。</p><h1 id="启用渐进式-JPEG-图片"><a href="#启用渐进式-JPEG-图片" class="headerlink" title="启用渐进式 JPEG 图片"></a>启用渐进式 JPEG 图片</h1><h2 id="渐进式图片转换"><a href="#渐进式图片转换" class="headerlink" title="渐进式图片转换"></a>渐进式图片转换</h2><p>因为上文提到的博客主题的更新，所以一进博客就要加载一张大背景图，如果还采用原来的线性加载的 jpg 图片的话会造成观感的不和谐，所以就将网站的图片进行了转换，使图片支持渐进式加载，这里的转换用到了 python 脚本</p><pre><code>from PIL import Image # pip3 install pilloworigin_file_path = &#39;./t.jpeg&#39;progressive_file_path = &#39;./o.jpeg&#39;original_image = Image.open(origin_file_path)original_image.convert(&#39;RGB&#39;)original_image.save(progressive_file_path, optimize=True, quality=100, progressive=True)</code></pre><h2 id="PNG-转-JPG"><a href="#PNG-转-JPG" class="headerlink" title="PNG 转 JPG"></a>PNG 转 JPG</h2><p>同样也使用了 python 脚本，这里顺便将图片也进行了渐进式 jpeg 的转换</p><pre><code>import osimport cv2import sysimport numpy as npfrom PIL import Image # pip3 install pillowpath = &quot;./&quot;print(path)for filename in os.listdir(path): if os.path.splitext(filename)[1] == &#39;.png&#39;: # print(filename) img = cv2.imread(path + filename) print(filename.replace(&quot;.png&quot;,&quot;.jpg&quot;)) newfilename = filename.replace(&quot;.png&quot;,&quot;.jpg&quot;) # cv2.imshow(&quot;Image&quot;,img) # cv2.waitKey(0) cv2.imwrite(path + newfilename,img) origin_file_path = path + newfilename progressive_file_path = path + newfilename original_image = Image.open(origin_file_path) original_image.convert(&#39;RGB&#39;) original_image.save(progressive_file_path, optimize=True, quality=100, progressive=True) os.remove(path+filename)% </code></pre><h1 id="iconfont-使用体验"><a href="#iconfont-使用体验" class="headerlink" title="iconfont 使用体验"></a>iconfont 使用体验</h1><p>写这套主题的时候，我有一点使用 icon 图标的需求，虽然说 <a href="https://fontawesome.com/" target="_blank" rel="noopener">fontawesome.com</a> 的图标品种十分丰富，但是似乎有些图标要使用的话要付费，对于我这种(穷)学生党来说有点难受，然后我想起了阿里巴巴开的 <a href="https://www.iconfont.cn/" target="_blank" rel="noopener">iconfont</a> 图标库，第一次使用之后就爱上了。</p><p><img src="https://cdn.lvcshu.info/img/20200419002.jpg" alt=""></p><p>它不仅提供常规的 icon 还提供了彩色的 icon，还是免费的(指没有商业使用的情况下)，爱了爱了，博客主题项目的图标就是使用的 iconfont 的图标</p><h1 id="更改-DNS-服务商"><a href="#更改-DNS-服务商" class="headerlink" title="更改 DNS 服务商"></a>更改 DNS 服务商</h1><p>之前使用的是 NS1 的免费 DNS 解析服务，虽然有分区域解析的功能，但是貌似效果不是很显著，趁着腾讯云的云解析正在搞活动买了一年的个人专业版(钱包 -￥36)来试用下，看了下可以区分境内境外解析，速度提升效果还可以</p><p><img src="https://cdn.lvcshu.info/img/20200420001.jpg" alt=""></p><p>EOF</p>',
          read: false
        }
      ],
      showPost: false,
      showRead: true,
      showUnread: true,
      nowPost: null
    };
  }
};
</script>
