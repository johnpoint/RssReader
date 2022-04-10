<template>
  <div>
    <h1 class="title">{{ post.Title }}</h1>
    <span>{{ post.Source }}</span> |
    <a :href="post.Url" target="_blank">{{ $t("post.link") }}</a>
    <div v-if="post.Content!==''" class="postcontext" style="margin: 15px" v-html="post.Content">
    </div>
    <div v-else class="postcontext" style="margin: 15px" v-html="post.Description">
    </div>
  </div>
</template>

<script>
export default {
  name: "Post",
  props: ['post', 'nextpost'],
}
</script>

<style scoped>

</style>