import React, { useState, useCallback, useEffect } from "react"
import { useRecoilValue, useSetRecoilState } from "recoil"
import entryDataAtom from "../../../recoil/entryData"
import SectionsRepresentation from "./SectionsRepresentation";
import { ReactComponent as ReplayIcon } from './replay.svg';
import styles from './EntryViewer.module.sass';
import { Tabs } from "../../UI";
import replayRequestModalOpenAtom from "../../../recoil/replayRequestModalOpen";

const enabledProtocolsForReplay = ["http"]

export enum TabsEnum {
    Request = 0,
    Response = 1
}

export const AutoRepresentation: React.FC<any> = ({ representation, color, openedTab = TabsEnum.Request, isDisplayReplay = false }) => {
    const entryData = useRecoilValue(entryDataAtom)
    const setIsOpenRequestModal = useSetRecoilState(replayRequestModalOpenAtom)
    const isReplayDisplayed = useCallback(() => {
        return enabledProtocolsForReplay.find(x => x === entryData.protocol.name) && isDisplayReplay
    }, [entryData.protocol.name, isDisplayReplay])

    // eslint-disable-next-line react-hooks/exhaustive-deps
    const TABS = [
        {
            tab: 'Request',
            badge: isReplayDisplayed() && <span title="Replay Request"><ReplayIcon fill={color} stroke={color} style={{ marginLeft: "10px", cursor: "pointer", height: "22px" }} onClick={() => setIsOpenRequestModal(true)} /></span>
        }
    ];
    const [currentTab, setCurrentTab] = useState(TABS[0].tab);

    const getOpenedTabIndex = useCallback(() => {
        const correntIndex = TABS.findIndex(current => current.tab === currentTab)
        return correntIndex > -1 ? correntIndex : 0
    }, [TABS, currentTab])

    useEffect(() => {
        if (openedTab) {
            setCurrentTab(TABS[openedTab].tab)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    // Don't fail even if `representation` is an empty string
    if (!representation) {
        return <React.Fragment></React.Fragment>;
    }

    const { request, response } = JSON.parse(representation);

    if (response) {
        TABS.push(
            {
                tab: 'Response',
                badge: null
            }
        );
    }

    return <div className={styles.Entry}>
        {<div className={styles.body}>
            <div className={styles.bodyHeader}>
                <Tabs tabs={TABS} currentTab={currentTab} color={color} onChange={setCurrentTab} leftAligned />
            </div>
            {getOpenedTabIndex() === TabsEnum.Request && <React.Fragment>
                <SectionsRepresentation data={request} color={color} requestRepresentation={request} />
            </React.Fragment>}
            {response && getOpenedTabIndex() === TabsEnum.Response && <React.Fragment>
                <SectionsRepresentation data={response} color={color} />
            </React.Fragment>}
        </div>}
    </div>;
}
