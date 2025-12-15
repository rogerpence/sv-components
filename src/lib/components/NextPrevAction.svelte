<script lang="ts">
	export interface Props {
		totalPages: number;
		pageNumber: number;
		navRoute: string;
	}

	let { totalPages, pageNumber, navRoute }: Props = $props();

	const nextPageNumber = $derived(pageNumber ? pageNumber + 1 : 0);
	const prevPageNumber = $derived(pageNumber ? pageNumber - 1 : 0);

	const isFirstPage = $derived(pageNumber == 1);
	const isLastPage = $derived(pageNumber == totalPages);

	// $inspect(isFirstPage, isLastPage, pageNumber);
</script>

<div class="page-navigator-container">
	<a
		aria-disabled={pageNumber === 1}
		class:disabled={pageNumber === 1}
		href="{navRoute}{prevPageNumber}"><i class="icon previous-icon"></i> Previous</a
	>
	Page {pageNumber} of {totalPages}
	<a
		aria-disabled={pageNumber === totalPages}
		class:disabled={pageNumber === totalPages}
		href="{navRoute}{nextPageNumber}">Next <i class="icon next-icon"></i></a
	>
</div>
