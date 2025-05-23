<script lang="ts" module>
  import { faCameraRetro } from '@fortawesome/sharp-duotone-light-svg-icons'
  export const page = {
    id: 'SnapshotApps',
    i: faCameraRetro,
    c1: 'burlywood',
    c2: 'darkgray',
    d1: 'burlywood',
    d2: 'silver',
  }
</script>

<script lang="ts">
  import { profile } from '../../api/profile.svelte'
  import { _ } from '../../includes/Translate.svelte'
  import { CardBody } from '@sveltestrap/sveltestrap'
  import Footer from '../../includes/Footer.svelte'
  import Header from '../../includes/Header.svelte'
  import Instances from '../../includes/Instances.svelte'
  import mysqlLogo from '../../assets/logos/mysql.png'
  import nvidiaLogo from '../../assets/logos/nvidia.png'
  import type { MySQLConfig, NvidiaConfig } from '../../api/notifiarrConfig'
  import { deepCopy } from '../../includes/util'
  import Instance, { type Form } from '../../includes/Instance.svelte'
  import InstanceHeader from '../../includes/InstanceHeader.svelte'
  import { InstanceFormValidator } from '../../includes/instanceFormValidator.svelte'
  import { nav } from '../../navigation/nav.svelte'

  const mysqlApp = {
    name: 'MySQL',
    id: page.id + '.MySQL',
    logo: mysqlLogo,
    hidden: ['deletes'],
    empty: {
      name: '',
      host: '',
      username: '',
      password: '',
      timeout: '10s',
      interval: '5m0s',
    },
    customValidator: (id: string, value: any) =>
      id.endsWith('.username')
        ? value === ''
          ? $_('phrases.UsernameMustNotBeEmpty')
          : ''
        : undefined,
    merge: (index: number, form: Form) => {
      const c = deepCopy($profile.config)
      c.snapshot!.mysql![index] = form as MySQLConfig
      return c
    },
  }

  const nvidiaApp = {
    name: 'Nvidia',
    id: page.id + '.Nvidia',
    logo: nvidiaLogo,
    hidden: ['deletes'],
    merge: (index: number, form: Form) => {
      const c = deepCopy($profile.config)
      c.snapshot!.nvidia = form as NvidiaConfig
      return c
    },
  }

  let iv = $derived({
    MySQL: new InstanceFormValidator($profile.config.snapshot?.mysql ?? [], mysqlApp),
    Nvidia: new InstanceFormValidator(
      [$profile.config.snapshot?.nvidia ?? {}],
      nvidiaApp,
    ),
  })

  async function submit() {
    const c = { ...$profile.config }
    c.snapshot!.mysql = iv.MySQL.instances as MySQLConfig[]
    c.snapshot!.nvidia = iv.Nvidia.instances[0]! as NvidiaConfig
    await profile.writeConfig(c)
  }

  $effect(() => {
    nav.formChanged = Object.values(iv).some(iv => iv.formChanged)
  })
</script>

<Header {page} />

<CardBody class="pt-0 mt-0">
  <Instances iv={iv.MySQL} />
  <InstanceHeader iv={iv.Nvidia} />
  <Instance
    bind:form={iv.Nvidia.instances[0]!}
    original={iv.Nvidia.original[0]!}
    app={nvidiaApp} />
</CardBody>

<Footer
  {submit}
  saveDisabled={Object.values(iv).every(iv => !iv.formChanged) ||
    Object.values(iv).some(iv => iv.invalid)} />
