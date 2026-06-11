import { computed, ref } from 'vue'
import type { WODType } from '@/features/wod/model/wodTypes'
import { SCORING_BY_TYPE } from '@/features/wod/model/wodTypes'

export function useWODTypes() {
  const selectedType = ref<WODType>('AMRAP')

  const scoringKind = computed(() => SCORING_BY_TYPE[selectedType.value])

  function setType(type: WODType) {
    selectedType.value = type
  }

  return {
    selectedType,
    scoringKind,
    setType,
  }
}
