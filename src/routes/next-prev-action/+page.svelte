<script lang="ts">
	import NextPrevAction from '$lib/components/NextPrevAction.svelte';
	import '$lib/css/next-prev-action/root.css';
	import '$lib/css/next-prev-action/utilities.css';
	import '$lib/css/next-prev-action/component.css';

	const totalPages = 21;
	const pageNumber = 1;
	const navRoute = '/technical-posts?pagenumber=';
</script>

<NextPrevAction {navRoute} {totalPages} {pageNumber}></NextPrevAction>

<style>
</style>
