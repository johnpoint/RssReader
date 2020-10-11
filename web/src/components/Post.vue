<template>
  <div>
    <h1 class="title">{{ post.Title }}</h1>
    <span>{{ post.Source }}</span> |
    <a :href="post.Url">{{ $t("post.link") }}</a>

    <b-card class="postcontext" style="margin: 15px" v-html="post.Description">
    </b-card>
  </div>
</template>

<script>
export default {
  name: "Post",
  props: ['post'],
  mounted() {
    console.log(this.post)
  }
}
</script>

<style scoped>

</style>