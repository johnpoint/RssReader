<template>
  <div class="home">
    <span>{{ info }}</span>
    <div id="list">
      <b-container class="tablist" fluid="true">
        <b-row>
          <b-col lg="0" @click="getPostList()" class="tab lefttab">{{
              $t("post.update")
            }}
          </b-col>
          <b-col lg="0"
                 @click="
          backToTop(5)
          " class="tab righttab">{{ $t("post.totop") }}
          </b-col>
        </b-row>


      </b-container>

      <label
          class="tab"
          :class="showUnread && !showRead ? 'select' : ''"
          @click="
          showUnread = true;
          showRead = false;
          nowData=post;
        "
      >{{ $t("post.unread") }}
        <a style="color: #42b983">{{ unreadpost }}</a>
      </label>
      |
      <label
          class="tab"
          :class="showRead && showUnread ? 'select' : ''"
          @click="
          showUnread = true;
          showRead = true;
        "
      >{{ $t("post.all") }}</label
      >
      |
      <label
          class="tab"
          :class="showRead && !showUnread ? 'select' : ''"
          @click="
          showUnread = false;
          showRead = true;
        "
      >{{ $t("post.read") }}</label
      >
      |
      <label
          class="tab"
          :class="!showRead && !showUnread ? 'select' : ''"
          @click="
          showUnread = false;
          showRead = false;
        "
      >{{ $t("post.readafter") }}</label
      >
      <transition-group name="flip-list" tag="div">
        <b-container v-for="(i,index) in cachePostData" :key="index" style="text-align: left">
          <b-row
              class="post"
              v-if="!showRead &&!showUnread&&i.readafter===true"
          >
            <b-col>
              <a style="font-size: small; color: rgba(0, 0, 0, 0.7)"
              >{{ i.Source }} >>
              </a>
              <a
                  style="font-size: large"
                  @click="
              toPost(cachePostData[index].ID)
            "
                  class="postlisttitle"
              >{{ i.Title }}
              </a>
            </b-col>
            <b-col cols="1" align="end">
            <b-icon-trash
                class="readbtn"
                style="float:right;text-align: right; margin: 5px;"
                @click="removeCache(index)"
            >saved
            </b-icon-trash>
            </b-col>
          </b-row>
        </b-container>
      </transition-group>

      <transition-group name="flip-list" tag="div">
        <b-container v-for="(i, index) in post" :key="index" style="text-align: left">
          <b-row
              class="post"
              :class="i.read ? 'read' : 'unread'"
              v-if="
            (showRead && i.read && !showUnread) ||
            (showUnread && !i.read) ||
            (showRead && showUnread)
          "
          >
            <b-col cols="0">
              <b-icon-check-square-fill
                  class="readbtn"
                  style="margin-top: 10px;color: #42b983"
                  v-if="i.read"
                  @click="change(index);"
              >read
              </b-icon-check-square-fill>
              <b-icon-check-square
                  class="readbtn"
                  style="margin-top: 10px;"
                  v-else
                  @click="change(index)"
              >unread
              </b-icon-check-square>
            </b-col>
            <b-col>
              <a
                  style="font-size: large;height: 22px"
                  @click="
              !post[index].read?change(index):'';
              toPost(post[index].ID);"
                  class="postlisttitle"
              >{{ i.Title }}
              </a>

              <a class="postdate">
                {{ i.date }}
              </a>
            </b-col>
            <b-col cols="0">
              <a class="postdate" style="font-size: small; color: rgba(0, 0, 0, 0.7);margin: 8px"
              >{{ i.Source }}
              </a>
            </b-col>
          </b-row>
        </b-container>
      </transition-group>
    </div>
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
    <span
        v-if="((showUnread && !showRead && unreadpost===0) || (!showRead&&!showUnread&&readafter.length===0)) && !showLoading"
        class="empty">{{ $t("post.empty") }}</span>
    <span v-else class="empty" style="color:rgba(0,0,0,0)">1</span>
  </div>
</template>

<script>
import axios from "axios";
import config from "@/config";
import router from "@/router";

export default {
  name: "Overview",
  components: {},
  beforeMount() {
    // check power
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true);
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"));
      if (window.localStorage.getItem("config") !== null) {
        this.$store.state.config = JSON.parse(window.localStorage.getItem("config"))
      } else {
        window.localStorage.setItem("config", JSON.stringify({"postnum": 50}))
      }
    } else {
      this.$store.commit("setStatus", false);
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    // get post list cache
    if (window.localStorage.getItem("posts") === null) {
      window.localStorage.setItem("posts", JSON.stringify([]));
    }
    this.post = JSON.parse(window.localStorage.getItem("posts"));
    this.$store.commit("setPostList", JSON.parse(window.localStorage.getItem("posts")))
    // something
    this.getPostList();
    this.updateCache();
  },
  beforeRouteLeave(to, from, next) {
    //console.log(to);
    //console.log(from);
    next();
  },
  methods: {
    removeCache: function (index) {
      // console.log("removeCache")
      window.localStorage.removeItem("post" + this.cachePostData[index].id);
      this.updateCache();
    },
    updateCache: function () {
      let cachePost = Object.keys(window.localStorage);
      this.cachePostList = [];
      this.cachePostData = [];
      for (let i1 = 0; i1 < cachePost.length; i1++) {
        let i = cachePost[i1];
        if (i.indexOf("post") !== -1 && i.indexOf("posts") === -1) {
          this.cachePostList.push(i.replace("post", ""))
          let data = JSON.parse(window.localStorage.getItem(i))
          data.id = parseInt(i.replace("post", ""))
          this.cachePostData.push(data)
        }
      }
      this.readafter = []
      for (let i of this.cachePostData) {
        if (i.readafter) {
          this.readafter.push(i)
        }
      }
      this.readafter.sort(function (m, n) {
        return -(m["time"] - n["time"])
      })
    },
    backToTop: function (rate) {
      var doc = document.body.scrollTop ? document.body : document.documentElement;
      var scrollTop = doc.scrollTop;

      var top = function () {
        scrollTop = scrollTop + (0 - scrollTop) / (rate || 2);

        // 临界判断，终止动画
        if (scrollTop <= 1) {
          doc.scrollTop = 0;
          return;
        }
        doc.scrollTop = scrollTop;
        // 动画gogogo!
        requestAnimationFrame(top);
      };
      top();
    },
    toPost: function (id) {
      this.$router.push("/post/" + id)
    },
    change: function (index) {
      // console.log("change")
      this.post[index].read ? this.unread(this.post[index].ID) : this.read(this.post[index].ID);
      this.post[index].read = !this.post[index].read;
    },
    read: function (id) {
      // console.log("read")
      if (id === 0) {
        id = this.nowshowpost.id
        if (this.nowshowpost.read === true) {
          return
        }
      }
      this.info = "";
      axios
          .post(
              config.apiAddress + "/api/post/read/" + id,
              null,
              {
                headers: {
                  Authorization: "Bearer " + this.$store.state.jwt,
                  Accept: "application/json",
                },
              }
          )
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                this.unreadpost -= 1;
                this.savePostListCache(JSON.stringify(this.post));
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                } else {
                  if (error.response.status === 401) {
                    window.localStorage.setItem("login", false)
                    router.push("/login")
                  }
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
    unread: function (id) {
      // console.log("unread")
      if (id === undefined) {
        id = this.nowshowpost.id
      }
      this.info = "";
      axios
          .post(
              config.apiAddress + "/api/post/unread/" + id,
              null,
              {
                headers: {
                  Authorization: "Bearer " + this.$store.state.jwt,
                  Accept: "application/json",
                },
              }
          )
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                this.unreadpost += 1;
                this.savePostListCache(JSON.stringify(this.post));
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                } else {
                  if (error.response.status === 401) {
                    window.localStorage.setItem("login", false)
                    router.push("/login")
                  }
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
    savePostListCache: function (data) {
      // console.log("savePostListCache")
      window.localStorage.setItem("posts", data);
      this.$store.commit("setPostList", JSON.parse(window.localStorage.getItem("posts")))
    },
    getReadList: function (postList) {
      // console.log("getReadList")
      this.unreadpost = 0;
      this.info = "";
      this.showLoading = true;
      axios
          .get(config.apiAddress + "/api/post/read", {
            headers: {
              Authorization: "Bearer " + this.$store.state.jwt,
              Accept: "application/json",
            },
          })
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                this.readPost = JSON.parse(response.data.message);
                this.post = [];
                this.unreadpost = 0;
                postList.forEach((item) => {
                  this.post.push({
                    ID: item.ID,
                    Title: item.Title,
                    Source: item.FeedTitle,
                    date: new Date(parseInt(item.Time) * 1000).format(
                        "yyyy-MM-dd hh:mm:ss"
                    ),
                    link: item.Link,
                    read: this.readPost.indexOf(item.ID) !== -1,
                    readAfter: false
                  });
                  this.readPost.indexOf(item.ID) === -1 ? this.unreadpost++ : null;
                });
                this.savePostListCache(JSON.stringify(this.post));
                this.showLoading = false;
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                } else {
                  if (error.response.status === 401) {
                    window.localStorage.setItem("login", false)
                    router.push("/login")
                  }
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
    getPostList: function () {
      // console.log("getPostList")
      this.info = "";
      let postNum = ""
      if (this.$store.state.config.postnum === undefined) {
        postNum = ""
      } else {
        postNum = this.$store.state.config.postnum
      }
      axios
          .get(config.apiAddress + "/api/post/" + postNum, {
            headers: {
              Authorization: "Bearer " + this.$store.state.jwt,
              Accept: "application/json",
            },
          })
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                let postList = JSON.parse(response.data.message);
                this.getReadList(postList);
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                  this.showLoading = false;
                } else {
                  if (error.response.status === 401) {
                    window.localStorage.setItem("login", false)
                    router.push("/login")
                  }
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
  },
  data() {
    return {
      post: [],//文章列表
      showRead: false,// tab控制开关
      showUnread: true,// tab控制开关
      readPost: [], // 已阅读文章列表
      showLoading: true,// 加载图标显示开关
      info: "",// 提示信息
      unreadpost: "-",// 未读文章计数
      cachePostList: [],// 缓存文章key
      cachePostData: [],// 缓存文章数据
      readafter: [],// 稍后阅读列表
      empty: false,// 判断文章列表是否为空
      nowshowpost: null,// 显示的文章
    };
  },
};
</script>
