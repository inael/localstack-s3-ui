import React, { useEffect, useState } from 'react'
import { WrapTable, ResourceRow, BreadCrums } from '../../Components'
import { withRouter } from 'react-router-dom'
import { joinPath } from '../../utils'

const apiURL = `http://localhost:8080/api/resource/`;

// TODO: need a linter. 
export const Resources = () => {

    const [state, setstate] = useState({
        path: "",
        resources: []
    })

    const fetchResources = async (resourcePath) => {
        try {
            const resourcesURL = resourcePath ? joinPath(apiURL, resourcePath) : apiURL
            console.log(resourcesURL)
            const res = await fetch(resourcesURL);
            const data = await res.json();
            console.log(data);
            setstate({
                path: data.path,
                name: data.name,
                type: data.type,
                resources: data.children || []
            })
        } catch (error) {
            console.log(error)
        }
    }

    useEffect(() => {
        fetchResources();
    }, [])

    const TableText = () => {
        if (state.name === "Root") {
            return <>
                <strong className="table-bucket-text">Buckets</strong>
                &nbsp;&nbsp;
                <strong className="table-bucket-nums">({state.resources.length})</strong>
            </>
        } else {
            return <strong className="table-bucket-text">{state.name}</strong>
        }
    }

    const TableHead = () => <>
        <BreadCrums path={state.path} fetchResources={fetchResources}/>
        <TableText />
    </>

    const TableBody = () =>
        state.resources.map((resource, index) =>
            <ResourceRow key={`bucketName-${index}`} resource={resource} fetchResources={fetchResources} />
        );

    return WrapTable(TableHead, TableBody)
}
