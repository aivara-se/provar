<script lang="ts">
  import { applicationStore } from '../../stores/application-store.svelte';
  import { projectStore } from '../../stores/project-store.svelte';
  import { Doctor } from '../../services/bindings';
  import type { bindings } from '../../../../wailsjs/go/models';
  import { X, CheckCircle, XCircle, AlertTriangle, Loader } from 'lucide-svelte';

  type DoctorCheck = bindings.DoctorCheck;

  type CheckState = 'pending' | 'running' | DoctorCheck;

  let checks = $state<CheckState[]>([]);
  let running = $state(false);

  const CHECK_LABELS = ['Playwright browser binaries', 'LLM API key', 'Project directory structure'];

  async function runChecks() {
    running = true;
    checks = CHECK_LABELS.map(() => 'running' as CheckState);
    try {
      const results = await Doctor.RunChecks(projectStore.path ?? '');
      checks = results;
    } catch (e) {
      console.error('DoctorModal: RunChecks failed:', e);
    } finally {
      running = false;
    }
  }

  $effect(() => {
    if (applicationStore.modalKind === 'doctor') {
      void runChecks();
    }
  });

  $effect(() => {
    if (applicationStore.modalKind !== 'doctor') return;
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') applicationStore.closeModal();
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  });

  function statusOf(c: CheckState): DoctorCheck['status'] | 'pending' | 'running' {
    if (c === 'pending') return 'pending';
    if (c === 'running') return 'running';
    return c.status;
  }

  function labelOf(c: CheckState, fallback: string): string {
    if (c === 'pending' || c === 'running') return fallback;
    return c.label;
  }

  function detailOf(c: CheckState): string {
    if (c === 'pending') return '';
    if (c === 'running') return 'Checking…';
    return c.detail;
  }
</script>

{#if applicationStore.modalKind === 'doctor'}
  <div
    class="fixed inset-0 z-[200] flex flex-col bg-[#0d1117] text-zinc-100"
    role="dialog"
    aria-label="Provar Doctor"
  >
    <!-- Header -->
    <header
      class="flex h-14 shrink-0 items-center justify-between border-b border-zinc-800 bg-[#161b22] px-6"
    >
      <div class="flex items-center gap-2.5">
      </div>
      <div class="flex items-center gap-3">
        <button
          type="button"
          disabled={running}
          class="rounded-md bg-indigo-600 px-4 py-1.5 text-xs font-medium text-white transition-colors hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
          onclick={runChecks}
        >
          {running ? 'Running…' : 'Re-run Checks'}
        </button>
        <button
          type="button"
          class="ml-2 rounded-md p-1.5 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-zinc-200"
          title="Close"
          onclick={() => applicationStore.closeModal()}
        >
          <X class="h-4 w-4" />
        </button>
      </div>
    </header>

    <!-- Body -->
    <div class="flex-1 overflow-y-auto p-8">
      <div class="mx-auto max-w-2xl space-y-4">
        <p class="mb-6 text-sm text-zinc-500">
          Checks your local environment for everything Provar needs to compile and run tests.
        </p>

        {#each checks as check, i (i)}
          {@const status = statusOf(check)}
          <div
            class="flex items-start gap-4 rounded-xl border bg-[#161b22] px-5 py-4 transition-colors {status ===
            'ok'
              ? 'border-emerald-800/60'
              : status === 'error'
                ? 'border-red-800/60'
                : status === 'warning'
                  ? 'border-amber-800/60'
                  : 'border-zinc-800'}"
          >
            <!-- Status icon -->
            <div class="mt-0.5 shrink-0">
              {#if status === 'running'}
                <Loader
                  class="h-5 w-5 animate-spin text-zinc-500"
                />
              {:else if status === 'ok'}
                <CheckCircle class="h-5 w-5 text-emerald-400" />
              {:else if status === 'error'}
                <XCircle class="h-5 w-5 text-red-400" />
              {:else if status === 'warning'}
                <AlertTriangle class="h-5 w-5 text-amber-400" />
              {:else}
                <div class="h-5 w-5 rounded-full border-2 border-zinc-700"></div>
              {/if}
            </div>

            <!-- Content -->
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium text-zinc-200">
                {labelOf(check, CHECK_LABELS[i])}
              </p>
              {#if detailOf(check)}
                <p
                  class="mt-1 font-mono text-xs {status === 'ok'
                    ? 'text-emerald-500'
                    : status === 'error'
                      ? 'text-red-400'
                      : status === 'warning'
                        ? 'text-amber-400'
                        : 'text-zinc-500'}"
                >
                  {detailOf(check)}
                </p>
              {/if}
            </div>

            <!-- Status badge -->
            {#if status !== 'pending' && status !== 'running'}
              <span
                class="shrink-0 rounded-full px-2 py-0.5 font-mono text-[10px] font-semibold uppercase tracking-widest {status ===
                'ok'
                  ? 'bg-emerald-900/40 text-emerald-400'
                  : status === 'error'
                    ? 'bg-red-900/40 text-red-400'
                    : 'bg-amber-900/40 text-amber-400'}"
              >
                {status}
              </span>
            {/if}
          </div>
        {/each}

        {#if checks.length === 0}
          <div
            class="flex items-center justify-center rounded-xl border border-zinc-800 bg-[#161b22] py-12"
          >
            <p class="text-sm text-zinc-600">No checks have run yet.</p>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
