<script lang="ts">
</script>

<header class="py-4 w-5/6 mx-auto text-center">
	<h1 class="font-black text-xl text-blue-500">koyo</h1>
	<p class="text-neutral-600">Website screenshot service api</p>
</header>

<hr class="w-1/2 mx-auto my-3 border-gray-100" />

<div>
	<div class="flex flex-col text-left w-3/5 mx-auto">
		<label for="website" class="text-neutral-700">Enter website url:</label>
		<div class="flex items-center">
			<input
				type="text"
				id="website"
				name="website"
				class="py-2 px-5 rounded-xl border-2 outline-none hover:border-blue-500 focus:border-blue-500 w-full"
				placeholder="Enter the website url..."
			/>

			<button
				class="inline-flex items-center text-sm bg-blue-400 hover:bg-blue-500 py-2 px-8 rounded-xl text-white duration-300 ml-2"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M9 17.25v1.007a3 3 0 01-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0115 18.257V17.25m6-12V15a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 15V5.25m18 0A2.25 2.25 0 0018.75 3H5.25A2.25 2.25 0 003 5.25m18 0V12a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 12V5.25"
					/>
				</svg>

				<span class="ml-1">Screenshot</span>
			</button>
		</div>
	</div>

	<div class="mt-12 border h-screen bg-gray-200 animate-pulse overflow-auto" />
</div>
