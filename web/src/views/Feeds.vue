<template>
  <div class="home">
    <span>{{ info }}</span><br v-if="info"/>
    <b-container class="tablist" fluid="true">
      <b-row>
        <b-col lg="0" class="tab lefttab" v-if="!addRss && !ioport"
               @click="addRss = true">{{ $t("feed.add") }}
        </b-col>
        <b-col lg="0" class="tab lefttab" v-else style="margin: 5px;" @click="addRss = false;ioport=false">{{
            $t("feed.cancel")
          }}
        </b-col>
        <b-col lg="0" class="tab righttab"
               v-if="!addRss && !ioport"
               @click="ioport=!ioport"
        >{{ $t("feed.import") }} / {{ $t("feed.export") }}
        </b-col>
      </b-row>
    </b-container>
    <div v-if="addRss" id="postinfo">
      <label>
        <input v-model="searchrss"/>
      </label>
      <b-button style="margin: 5px" @click="searchRss()">{{ $t("feed.search") }}</b-button>
      <br>
      <b-container>
        <b-row>
          <b-col v-for="(i, index) in search" :key="index" style="text-align: left">
            <b-card>
              <p>{{ i.title }}</p><span style="font-size: small">{{ i.link }}</span>
              <hr>
              <span>{{ $t("feed.subscriber") }} {{ i.subscriber }}</span>
              <b-icon-plus-square
                  style="margin: 5px;float: right"
                  @click="addSub(index)"
              ></b-icon-plus-square>
            </b-card>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <div id="ioport" v-if="ioport" style="max-width: 480px;margin: auto">
      <b-card style="text-align: left">
        <label>{{ $t("feed.import") }}</label>
        <b-form-file
            v-model="opml"
            :state="Boolean(opml)"
            placeholder="Choose a file or drop it here..."
            drop-placeholder="Drop file here..."
        ></b-form-file>
        <b-button size="sm" style="float: right;margin: 5px" @click="opml=null" variant="outline-primary">{{
            $t("feed.clear")
          }}
        </b-button>
        <div class="mt-3">Selected file: {{ opml ? opml.name : '' }}</div>
        <b-button size="sm" style="margin: 5px" @click="uploadopml()" variant="outline-success">{{
            $t("feed.upload")
          }}
        </b-button>
      </b-card>
      <b-card style="text-align: left">
        <label>{{ $t("feed.export") }}</label><br>
        <b-button size="sm" @click="exportopml()" style="margin: 5px" variant="outline-primary">{{
            $t("feed.download")
          }}
        </b-button>
      </b-card>
    </div>
    <label v-if="!addRss && !ioport" style="margin: 5px;font-size: larger">{{ $t("feed.subscribed") }}</label>
    <div v-if="!addRss && !ioport" id="list">
      <b-container v-for="(i, index) in rss" :key="index" style="text-align: left">
        <b-row class="post">
          <b-col cols="0"><a
              style="font-size: small;color: rgba(0,0,0,.7);margin: 5px"
          >{{ i.unread }}
          </a></b-col>
          <b-col col="11">
            <a style="font-size: large;color: rgba(0, 0, 0, 0.7)">{{ i.title }} </a><span
              style="font-size: small">{{ i.link }} </span>
            <b-icon-exclamation-circle v-if="i.status>0" v-b-tooltip.hover :title='$t("feed.geterror")'
                                                                                      style="color: red"></b-icon-exclamation-circle>
            <b-icon-exclamation-circle v-if="i.status===-1" v-b-tooltip.hover :title='$t("feed.getstop")'
                                       style="color: red"></b-icon-exclamation-circle>
          </b-col>
          <b-col cols="4" align="end">
            <b-icon-check-square-fill
                style="color: rgb(69,123,48);margin: 5px"
                v-if="i.unread === 0"
            >read
            </b-icon-check-square-fill
            >
            <b-icon-check-square
                v-else
                style="margin: 5px"
                @click="i.unread = 0;readFeed(index)"
            >unread
            </b-icon-check-square
            >
            <b-icon-x
                v-if="delRss && delRssIndex === index"
                style="margin: 5px"
                @click="delRss = false"
            ></b-icon-x>
            <b-icon-check
                v-if="delRss && delRssIndex === index"
                style="margin: 5px"
                @click="removeRss(index);delRss = false"
            ></b-icon-check>
            <b-icon-trash
                v-else
                style="margin: 5px"
                @click="
              delRss = true;
              delRssIndex = index;
            "
            ></b-icon-trash>
          </b-col>
        </b-row>
      </b-container>
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

export default {
  name: "Overview",
  components: {},
  beforeMount() {
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
    this.getData()
    this.getRss()
  },
  methods: {
    "saveData": function () {
      window.localStorage.setItem("feeds", JSON.stringify(this.rss))
    },
    "getData": function () {
      if (window.localStorage.getItem("feeds") === null) {
        window.localStorage.setItem("feeds", JSON.stringify([]))
      }
      this.rss = JSON.parse(window.localStorage.getItem("feeds"))
    },
    "readFeed": function (index) {
      axios.post(config.apiAddress + "/api/feed/read/" + this.rss[index].id, null, {
        "headers": {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(
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
          })
    },
    "change": function (index) {
      this.post[index].read = !this.post[index].read;
    },
    "addSub": function (index) {
      this.info = ""
      this.showLoading = true
      axios.post(config.apiAddress + "/api/feed/subscribe/" + this.search[index].id, null, {
        "headers": {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(
          (response) => {
            if (response.data.code === 200) {
              this.getRss()
              this.addRss = false
              this.searchrss = ""
              this.search = []
            } else {
              this.info = response.data.message
            }
            this.showLoading = false
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
          })
    },
    "removeRss": function (index) {
      this.showLoading = true
      axios.post(config.apiAddress + "/api/feed/unsubscribe/" + this.rss[index].id, null, {
        "headers": {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(() => {
        this.getRss()
      })
      this.showLoading = false
    },
    "searchRss": function () {
      this.showLoading = true
      this.info = ""
      axios.post(config.apiAddress + "/api/feed/search", {
        "Url": this.searchrss
      }, {
        "headers": {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(
          (response) => {
            if (response.data.code === 200) {
              let data = JSON.parse(response.data.message)
              this.search = []
              this.search.push({
                "id": data.ID,
                "title": data.Title,
                "link": data.Url,
                "subscriber": data.Subscriber
              })
              this.showLoading = false
            } else {
              this.info = response.data.message
              this.showLoading = false
            }
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
          })
    },
    "getRss": function () {
      this.info = ""
      axios.get(config.apiAddress + "/api/feed/list", {
        "headers": {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(
          (response) => {
            if (response.data.code !== 200) {
              this.info = response.data.message
              return
            }
            let postnum = ""
            if (this.$store.state.config.postnum === undefined) {
              postnum = ""
            } else {
              postnum = this.$store.state.config.postnum
            }
            this.rsslist = JSON.parse(response.data.message);
            axios.get(config.apiAddress + "/api/post/" + postnum, {
              "headers": {
                'Authorization': "Bearer " + this.$store.state.jwt,
                'Accept': 'application/json'
              }
            }).then(
                (response) => {
                  if (response.data.code !== 200) {
                    this.info = response.data.message
                    return
                  }
                  this.postList = JSON.parse(response.data.message)
                  axios.get(config.apiAddress + "/api/post/read", {
                    "headers": {
                      'Authorization': "Bearer " + this.$store.state.jwt,
                      'Accept': 'application/json'
                    }
                  }).then(
                      (response) => {
                        if (response.data.code !== 200) {
                          this.info = response.data.message
                          return
                        }
                        this.rss = []
                        this.readPost = JSON.parse(response.data.message)
                        this.rsslist.forEach(item => {
                          let unread = 0
                          this.postList.forEach(post => {
                            if (post.Feed === item.ID && this.readPost.indexOf(post.ID) === -1) {
                              unread++
                            }
                          })

                          this.rss.push({
                            "id": item.ID,
                            "title": item.Title,
                            "link": item.Url,
                            "unread": unread,
                            "status": item.Status
                          })
                        })
                        this.saveData()
                        this.showLoading = false
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
                      })
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
                })
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
          })
    },
    "exportopml": function () {
      this.showLoading = true;
      axios.get(config.apiAddress + "/api/user/opml", {
        "headers": {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json',
        }
      }).then(
          (response) => {
            if (response.data.code !== 200) {
              this.info = response.data.message;
              return
            }
            var pom = document.createElement('a');
            pom.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(response.data.message));
            pom.setAttribute('download', "export.xml");

            if (document.createEvent) {
              var event = document.createEvent('MouseEvents');
              event.initEvent('click', true, true);
              pom.dispatchEvent(event);
            } else {
              pom.click();
            }
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
          })
    },
    "uploadopml": function () {
      if (this.opml !== null) {
        var formData = new FormData();
        formData.append("opml", this.opml);
        axios.post(config.apiAddress + "/api/user/opml",
            formData
            , {
              "headers": {
                'Authorization': "Bearer " + this.$store.state.jwt,
                'Accept': 'application/json',
                'Content-Type': 'multipart/form-data'
              }
            }).then(
            (response) => {
              if (response.data.code !== 200) {
                this.info = response.data.message
                return
              }
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
            })
      }
    }
  }
  ,
  data() {
    return {
      rss: [],
      info: "",
      addRss: false,
      delRss: false,
      delRssIndex: 0,
      rsslist: null,
      readPost: null,
      searchrss: "",
      unread: 0,
      search: [],
      showLoading: true,
      dismissSecs: 10,
      dismissCountDown: 0,
      showDismissibleAlert: false,
      postList: [],
      ioport: false,
      opml: null
    };
  }
}
;
</script>
