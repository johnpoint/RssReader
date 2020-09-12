<template>
  <div class="home">
    <check />
    <div v-if="!addRss" id="list">
      <label class="tab" @click="addRss = true">Add</label>
      <div v-for="(i, index) in rss" :key="index" style="text-align: left">
        <div class="post">
          <a style="font-size: large">{{ i.title }} </a>
          <span style="font-size: small">{{ i.date }}</span>
          <b-icon-x
            v-if="delRss && delRssIndex == index"
            style="float: right;margin: 5px"
            @click="delRss = false"
          ></b-icon-x>
          <b-icon-check
            v-if="delRss && delRssIndex == index"
            style="float: right;margin: 5px"
            @click="removeRss(index)"
          ></b-icon-check>
          <b-icon-trash
            v-else
            style="float: right;margin: 5px;"
            @click="
              delRss = true;
              delRssIndex = index;
            "
          ></b-icon-trash>
          <b-icon-check-square-fill
            style="float: right;margin: 5px;color: rgb(69,123,48)"
            v-if="i.unread == 0"
            >read</b-icon-check-square-fill
          >
          <b-icon-check-square
            style="float: right;margin: 5px"
            v-else
            @click="i.unread = 0"
            >unread</b-icon-check-square
          >
          <a
            style="font-size: small;color: rgba(0,0,0,.7);float: right;margin: 5px"
          >
            {{ i.unread }}
          </a>
        </div>
      </div>
    </div>
    <div v-else id="postinfo">
      <label class="tab" style="margin: 5px;float: left" @click="addRss = false"
        >Back</label
      >
      <input v-model="searchrss" />
      <b-button style="margin: 5px" @click="getRss()">ok</b-button>
      <div v-for="(i, index) in search" :key="index" style="text-align: left">
        <div class="post">
          <a>{{ i.title }} </a>
          <span style="font-size: small"> {{ i.link }}</span>
          <b-icon-plus-square
            style="margin: 5px;float: right"
            @click="addSub(index)"
          ></b-icon-plus-square>
        </div>
      </div>
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
      this.post[index].read = this.post[index].read ? false : true;
    },
    addSub: function(index) {
      console.log(index);
    },
    removeRss: function(index) {
      console.log(index);
    },
    getRss: function() {
      console.log(this.searchrss);
    }
  },
  data() {
    return {
      rss: [
        {
          title: "johnpoint's blog",
          date: "2020-06-13 13:00",
          link: "https://blog.lvcshu.com",
          unread: 3
        },
        {
          title: "johnpoint's blog",
          date: "2020-06-13 13:00",
          link: "https://blog.lvcshu.com",
          unread: 3
        }
      ],
      addRss: false,
      delRss: false,
      delRssIndex: 0,
      searchrss: "",
      search: [
        {
          title: "johnpoint's blog",
          date: "2020-06-13 13:00",
          link: "https://blog.lvcshu.com"
        }
      ]
    };
  }
};
</script>
