<script lang="ts">
	export let websiteUrl: string;
	export let optionsImageType: string;
	export let imageUrl: string;
</script>

<div>
	<a
		download={`${websiteUrl.replace('https://', '').replace('http://', '')}.${optionsImageType}`}
		href={imageUrl}
		class="inline-flex items-center text-sm py-2 px-6 rounded-xl bg-blue-400 hover:bg-blue-500 text-white duration-300"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="w-5 h-5"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"
			/>
		</svg>

		<span class="ml-1">Download</span>
	</a>
</div>
