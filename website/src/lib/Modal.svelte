<script lang="ts">
	import { Dialog, TransitionChild, Transition, DialogOverlay } from '@rgossiaux/svelte-headlessui';
	export let isOpen: boolean;
	export let className: string;
	export let closeModal: () => void;
</script>

<Transition appear show={isOpen}>
	<Dialog as="div" class="fixed inset-0 z-10 overflow-y-auto" on:close={closeModal}>
		<div class="min-h-screen px-4 text-center flex items-center justify-center">
			<TransitionChild
				enter="ease-out duration-300"
				enterFrom="opacity-0"
				enterTo="opacity-100"
				leave="ease-in duration-200"
				leaveFrom="opacity-100"
				leaveTo="opacity-0"
			>
				<DialogOverlay class="fixed inset-0 bg-black bg-opacity-25" />
			</TransitionChild>

			<!-- This element is to trick the browser into centering the modal contents. -->
			<!-- <span class="inline-block h-screen align-middle" aria-hidden="true"> &#8203; </span> -->

			<div
				class={`inline-block w-full p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl ${className}`}
			>
				<slot />
			</div>
		</div>
	</Dialog>
</Transition>
