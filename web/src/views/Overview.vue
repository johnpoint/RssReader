<template>
  <div class="home">
    <check/>
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
          >read
          </b-icon-check-square-fill
          >
          <b-icon-check-square
              style="float: right;margin: 5px"
              v-else
              @click="change(index)"
          >unread
          </b-icon-check-square
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
        >read
        </b-icon-check-square-fill
        >
        <label v-if="post[nowPost].read">已读</label>
        <b-icon-check-square
            style="float: right;margin: 5px"
            v-if="!post[nowPost].read"
            @click="change(nowPost)"
        >unread
        </b-icon-check-square
        >
        <label v-if="!post[nowPost].read">未读</label>
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
import axios from "axios"
import config from "@/config";

export default {
  name: "Overview",
  components: {
    check
  },
  beforeMount() {
    this.getPostList()
  },
  methods: {
    change: function (index) {
      this.post[index].read ? this.unread(index) : this.read(index);
      this.post[index].read = this.post[index].read ? false : true;
    },
    read: function (index) {
      axios.post(config.apiAddress + "/api/post/read/" + this.post[index].id, null,{
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        console.log(response.data)
      })
    },
    unread: function (index) {
      console.log(index);
    },
    getReadList: function () {
      axios.get(config.apiAddress + "/api/post/read", {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        this.readPost = JSON.parse(response.data.message)
        this.postList.forEach(item => {
          this.post.push({
            id: item.ID,
            title: item.Title,
            date: item.Time,
            read: this.readPost.indexOf(item.ID) == -1 ? false : true
          })
        })
      })
    },
    getPostList: function () {
      axios.get(config.apiAddress + "/api/post/", {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        this.postList = JSON.parse(response.data.message)
        this.getReadList()
      })
    },
    getPostContent: function () {
      return;
    }
  },
  data() {
    return {
      post: [],
      showPost: false,
      showRead: true,
      showUnread: true,
      nowPost: null,
      readPost: [],
      postList: null
    };
  }
};
</script>
