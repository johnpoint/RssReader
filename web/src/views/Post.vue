<template>
  <div id="postinfo">
    <b-container class="tablist" fluid="true">
      <b-row>
        <b-col
            lg="0"
            @click="back"
            class="tab lefttab"
        >{{ $t("post.back") }}
        </b-col
        >
        <b-col
            lg="0"
            @click="readAfter(nowshowpost);$router.push('/posts')"
            class="tab righttab"
        >{{ $t("post.setunread") }}
        </b-col
        >
      </b-row>
    </b-container>
    <div class="postbox">
      <b-skeleton-wrapper>
        <post v-if="nowshowpost!==null" :post="nowshowpost" :nextpost="nextpost"/>

        <center v-else>
          <b-skeleton style="text-align: center" type="h1" class="title" width="50%"></b-skeleton>
          <br/>
          <b-skeleton style="a" width="30%"></b-skeleton>
          <div class="postcontext" style="margin: auto;">
            <b-skeleton width="85%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="55%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="70%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="85%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="55%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="70%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="55%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="70%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="85%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="85%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="55%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="55%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="70%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="55%" style="margin: auto;margin-top: 4px"></b-skeleton>
            <b-skeleton width="70%" style="margin: auto;margin-top: 4px"></b-skeleton>
          </div>
        </center>
      </b-skeleton-wrapper>

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
      nextpost: null,
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
    } else {
      this.$store.commit("setStatus", false);
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    if (window.localStorage.getItem("posts") === null) {
      window.localStorage.setItem("posts", JSON.stringify([]));
    }
    this.post = JSON.parse(window.localStorage.getItem("posts"));
    this.$store.commit("setPostList", JSON.parse(window.localStorage.getItem("posts")))
  },
  mounted() {
    this.getPostContent();
  },
  beforeRouteLeave(to, from, next) {
    //console.log(to);
    //console.log(from);
    next();
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
          data.ID = parseInt(i.replace("post", ""))
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
      console.log(data)
      let now = JSON.parse(window.localStorage.getItem("post" + data._id))
      now.readafter = true
      now.time = Date.parse(new Date())
      window.localStorage.setItem("post" + data._id, JSON.stringify(now))
      this.showPost = false
      this.updateCache()
    },
    back: function () {
      this.$router.push("/posts")
    },
    getPostContent: function () {
      let id = this.$route.params.id
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
                let newPostCache = response.data.data
                for (let i of this.post) {
                  if (i.ID === newPostCache._id) {
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
                console.log(newPostCache)
                this.nowshowpost = newPostCache
                // console.log(newPostCache)
                this.updateCache();
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
  }
}
</script>

<style scoped>

</style>