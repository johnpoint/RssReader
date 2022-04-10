<template>
  <div>
    <h1 class="title">{{ post.title }}</h1>
    <span>{{ post.Source }}</span> |
    <a :href="post.url" target="_blank">{{ $t("post.link") }}</a>
    <div v-if="post.content!==''" class="postcontext" style="margin: 15px" v-html="post.content">
    </div>
    <div v-else class="postcontext" style="margin: 15px" v-html="post.description">
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