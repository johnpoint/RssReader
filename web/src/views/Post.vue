<template>
  <div id="postinfo">
    <div class="tablist">
      <label
          @click="back"
          class="tab lefttab"
      >{{ $t("post.back") }}</label
      >
      <label
          @click="readAfter(nowshowpost);"
          class="tab righttab"
      >{{ $t("post.setunread") }}</label
      >
    </div>
    <div class="postbox">
      <post v-if="nowshowpost!==null" :post="nowshowpost"/>
    </div>
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import config from "@/config";
import router from "@/router";
import post from "@/components/Post"


export default {
  name: "Post",
  components: {
    post
  },
  data() {
    return {
      nowshowpost: null,// 显示的文章
      loadingShowPost: {
        Title: "Loading Title",
        Url: "/",
        Source: "Loading",
        Description: "Loading"
      },
      showLoading: false,
    }
  },
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
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    if (window.localStorage.getItem("posts") === null) {
      window.localStorage.setItem("posts", JSON.stringify([]));
    }
    this.post = JSON.parse(window.localStorage.getItem("posts"));
    this.$store.state.postList = JSON.parse(window.localStorage.getItem("posts"));
    this.getPostContent();
  },
  methods: {
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
    readAfter: function (data) {
      // console.log("readAfter")
      if (data.read !== true) {
        let j = 0;
        for (let i of this.post) {
          if (i.ID === data.ID) {
            this.change(j)
            break
          }
          j++;
        }
      }
      let now = JSON.parse(window.localStorage.getItem("post" + data.ID))
      now.readafter = true
      now.time = Date.parse(new Date())
      window.localStorage.setItem("post" + data.ID, JSON.stringify(now))
      this.showPost = false
      this.updateCache()
    },
    back: function () {
      this.$router.push("/posts")
    },
    getPostContent: function () {
      console.log("getPostContent")
      let id = this.$route.params.id
      console.log(id)
      if (window.localStorage.getItem("post" + id) !== null) {
        this.nowshowpost = JSON.parse(window.localStorage.getItem(
            "post" + id
        ));
        this.showLoading = false;
        return;
      }
      this.info = "";
      this.showLoading = true;
      axios
          .get(config.apiAddress + "/api/post/content/" + id, {
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
                let newPostCache = JSON.parse(response.data.message)
                for (let i of this.post) {
                  if (i.ID === newPostCache.ID) {
                    newPostCache.Source = i.Source
                    newPostCache.read = i.read
                    break
                  }
                }
                newPostCache.readafter = false
                window.localStorage.setItem(
                    "post" + id,
                    JSON.stringify(newPostCache)
                );
                this.showLoading = false;
                this.nowshowpost = newPostCache
                // console.log(newPostCache)
                this.updateCache();
              },
              (error) => {
                let errText
                if (error.response == undefined) {
                  errText = "Unable to connect to server";
                } else {
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
  }
}
</script>

<style scoped>

</style>